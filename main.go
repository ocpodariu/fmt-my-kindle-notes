package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/html"
)

var (
	outputFilename   = flag.String("out", "notes.html", "Name of the output HTML file")
	templateFilename = flag.String("template", "output.tpl", "Template for the formatted notebook")
	verboseMode      = flag.Bool("v", false, "Enable verbose mode")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: fmt-my-kindle-notes [OPTIONS] NOTEBOOK")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.SetOutput(ioutil.Discard)
	if *verboseMode {
		log.SetFlags(0)
		log.SetPrefix("> ")
		log.SetOutput(os.Stdout)
	}

	notebookPath := args[0]
	f, err := os.Open(notebookPath)
	if err != nil {
		fmt.Printf("open notebook: %v", err)
		os.Exit(2)
	}

	log.Printf("Parsing notebook '%s'", notebookPath)
	htmlRoot, err := html.Parse(f)
	if err != nil {
		fmt.Printf("parse HTML file: %v", err)
		os.Exit(2)
	}

	log.Println("Extracting highlights and notes")
	var notebook Notebook
	if err = notebook.Parse(htmlRoot); err != nil {
		fmt.Printf("parse notebook: %v", err)
		os.Exit(2)
	}

	tpl, err := template.ParseFiles(*templateFilename)
	if err != nil {
		fmt.Printf("parse output template: %v", err)
		os.Exit(2)
	}
	g, err := os.Create(*outputFilename)
	if err != nil {
		fmt.Printf("create output file: %v", err)
		os.Exit(2)
	}
	log.Printf("Writing formatted notebook to '%s'", *outputFilename)
	err = tpl.Execute(g, notebook)
	if err != nil {
		fmt.Printf("render template: %v", err)
		os.Exit(2)
	}
}
