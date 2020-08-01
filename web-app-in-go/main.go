package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

const extension = ".txt"

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "delete.html"))

type page struct {
	Title string
	Body  []byte
}

func (p *page) save() error {
	filename := p.Title + extension
	return ioutil.WriteFile(filename, p.Body, 0600) //WriteFile is a standard library function that writes a byte slice to a file
}

func deletePage(title string) error {
	filename := title + extension
	return os.Remove(filename)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/delete/"):]
	err := deletePage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, _ := template.ParseFiles("delete.html")
	t.Execute(w, r)
}

func renderTemplate(w http.ResponseWriter, templ string, p *page) {
	err := templates.ExecuteTemplate(w, templ+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func loadPage(title string) (*page, error) {
	filename := title + extension
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/delete/", deleteHandler)
	fmt.Print("Listenning on 8080 ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
