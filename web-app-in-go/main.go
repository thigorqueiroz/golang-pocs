package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

const extension = ".txt"

var templates = template.Must(template.ParseGlob("html/*"))

var validPath = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9]+)$")

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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	t := r.URL.Path[len("/delete/"):]
	err := deletePage(t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ, _ := template.ParseFiles("html/delete.html")
	templ.Execute(w, r)
}

func renderTemplate(w http.ResponseWriter, templ string, p *page) {
	err := templates.ExecuteTemplate(w, templ+".html", p)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &page{Title: title, Body: []byte(body)}
	errr := p.save()
	if errr != nil {
		log.Println(errr)
		http.Error(w, errr.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func loadPage(title string) (*page, error) {
	filename := title + extension
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))
	fmt.Print("Listenning on 8080 ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
