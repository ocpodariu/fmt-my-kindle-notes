## fmt-my-kindle-notes

Parse Kindle notes exported as HTML and format them however you want.

You can write your own custom HTML templates using Go's templating syntax (see [Golang text/template](https://golang.org/pkg/text/template/) or [Golang Template Cheatsheet](https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet)).

I've included a sample notebook (`sample/kindle-notebook.html`) to make it easier to test this tool.


### How to install
Get [Git](https://git-scm.com/) and [Go](https://golang.org/doc/install) in case you don't have them already, clone the repository and build the app.
```
git clone https://github.com/ocpodariu/fmt-my-kindle-notes.git
cd fmt-my-kindle-notes
go build
```

### How to use
```
./fmt-my-kindle-notes [OPTIONS] NOTEBOOK
```

#### Flags
`-out` (default: "notes.html") - Change the name of the output HTML file:
```
./fmt-my-kindle-notes -out book_notes.html NOTEBOOK
```

`-template` (default: "output.tpl") - Use a custom template for the output HTML file:
```
./fmt-my-kindle-notes -template awesome.tpl NOTEBOOK
```

`-v` (default: disabled) - Enable verbose mode to view debug messages
```
./fmt-my-kindle-notes -v NOTEBOOK
```

