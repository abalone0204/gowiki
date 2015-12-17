package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("view.html", "edit.html"))

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func viewHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	p, err := load(title)
	if err != nil {
		http.Redirect(w, req, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
	body := req.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Page
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, err
}

func main() {
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
