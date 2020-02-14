package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

const (
	SectionClass   = "sectionHeading"
	HighlightClass = "noteHeading"
	NoteTextClass  = "noteText"
	AuthorClass    = "authors"
	TitleClass     = "bookTitle"
)

type Highlight struct {
	Color    string
	Page     int
	Location int
	Text     string
	Note     string
}

type Section struct {
	Title      string
	Highlights []*Highlight
}

type Notebook struct {
	Title    string
	Author   string
	Sections []*Section
}

// Parse extracts sections, highlights and notes from a Kindle notebook
// exported as HTML
func (nb *Notebook) Parse(htmlNotebook *html.Node) error {
	// All highlights and notes are in a div element with class="bodyContainer".
	// Find this element and then iterate its children.
	bc := findFirstNodeByClass(htmlNotebook, "bodyContainer")
	if bc == nil {
		return ErrNoBodyContainer
	}

	var (
		crtSection   *Section
		crtHighlight *Highlight

		// Indicates that the next div element with class NoteTextClass
		// contains the text of a note, not the text of a highlight.
		//
		// NB: A note always appears immediately after a highlight.
		insideNote bool
	)

	for c := bc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}

		switch getClassAttr(c) {
		case SectionClass:
			// A new section starts
			// Do this check to avoid adding a nil Section on the first iteration
			if crtSection != nil {
				// Add previous section to notebook
				crtSection.Highlights = append(crtSection.Highlights, crtHighlight)
				nb.Sections = append(nb.Sections, crtSection)

				crtHighlight = nil // Reset
			}
			crtSection = &Section{Title: strings.TrimSpace(c.FirstChild.Data)}

		case HighlightClass:
			// A new Highlight or an attached Note starts
			highlightText := strings.TrimLeftFunc(c.FirstChild.Data, unicode.IsSpace)
			if strings.HasPrefix(highlightText, "Highlight") {
				// Do this check to avoid adding a nil Highlight on first iteration
				if crtHighlight != nil {
					crtSection.Highlights = append(crtSection.Highlights, crtHighlight)
				}
				crtHighlight = parseHighlightNode(c)
			} else if strings.HasPrefix(highlightText, "Note") {
				insideNote = true
			}

		case NoteTextClass:
			// The text of a highlight or note
			text := parseTextNode(c)
			if insideNote {
				crtHighlight.Note = text
				insideNote = false // Reset
			} else {
				crtHighlight.Text = text
			}

		case AuthorClass:
			nb.Author = strings.TrimSpace(c.FirstChild.Data)

		case TitleClass:
			nb.Title = strings.TrimSpace(c.FirstChild.Data)
		}
	}
	// Append last highlight and section
	crtSection.Highlights = append(crtSection.Highlights, crtHighlight)
	nb.Sections = append(nb.Sections, crtSection)

	return nil
}

func (nb *Notebook) Print() {
	fmt.Printf("'%s' by %s\n", nb.Title, nb.Author)
	for i, s := range nb.Sections {
		fmt.Printf("Section %02d: %s\n", i+1, s.Title)
		for j, h := range s.Highlights {
			fmt.Printf("Highlight %02d-%d\n", i+1, j+1)
			fmt.Printf("\tText: %s\n", h.Text)
			fmt.Printf("\tNote: %s\n", h.Note)
		}
		fmt.Println()
	}
}

var ErrNoBodyContainer = errors.New(`expected element with class="bodyContainer"`)

var (
	reExtractLocation = regexp.MustCompile(`Location (\d+)`)
	reExtractPage     = regexp.MustCompile(`Page (\d+)`)
)

// parseHighlightNode parses a Highlight node to determine the page, location
// and color of the highlight.
//
// A Highlight node looks like this (contains only location):
// <div class="noteHeading">
//     Highlight (<span class="highlight_yellow">yellow</span>) -  Location 2364
// </div>
//
// Or like this (contains both page and location):
// <div class="noteHeading">
//     Highlight(<span class="highlight_yellow">yellow</span>) - 10. Scientists as Explorers of the Universe  > Page 115 · Location 954
// </div>
func parseHighlightNode(n *html.Node) *Highlight {
	// The Node structure for a Highlight node has three children:
	// 1. TextNode with the first part of the text:
	//     "Highlight("
	// 2. ElementNode with the highlight color:
	//     <span class="highlight_yellow">yellow</span>
	// 3. TextNode with the second part of the text:
	//     ") -  Location 2364"

	var h Highlight
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			matches := reExtractLocation.FindStringSubmatch(c.Data)
			if matches != nil {
				// It's safe to skip the error check because the regular expression matches only numbers, and location numbers definitely fit in an int64
				loc, _ := strconv.Atoi(matches[1])
				h.Location = loc
			}

			matches = reExtractPage.FindStringSubmatch(c.Data)
			if matches != nil {
				// It's safe to skip the error check because the regular expression matches only numbers, and page numbers definitely fit in an int64
				page, _ := strconv.Atoi(matches[1])
				h.Page = page
			}

		case html.ElementNode:
			// Extract the highlight color from the class attribute
			colorClass := getClassAttr(c)
			idx := strings.IndexRune(colorClass, '_')
			if idx != -1 {
				h.Color = colorClass[idx+1:]
			} else {
				// HTML might be malformed.
				// Expected to find highlight_<color_name>
			}
		}
	}

	return &h
}

// parseTextNode parses Note nodes and returns the their content.
//
// A Note node looks like this:
// <div class="noteText">
//   Entrepreneurship is enhanced by performing lots of quick, easily performed experiments.
// </div>
func parseTextNode(n *html.Node) string {
	if getClassAttr(n) != NoteTextClass {
		return ""
	}
	return strings.TrimSpace(n.FirstChild.Data)
}

// findFirstNodeByClass does a depth-first search and returns
// the first node that has the given class attribute
func findFirstNodeByClass(n *html.Node, className string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == className {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nc := findFirstNodeByClass(c, className); nc != nil {
			return nc
		}
	}

	return nil
}

// getClassAttr returns the value of a node's "class" attribute.
// If the node doesn't have a "class" attribute, then it return an empty string.
func getClassAttr(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "class" {
			return a.Val
		}
	}
	return ""
}
