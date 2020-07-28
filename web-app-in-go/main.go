package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const extension = ".txt"

type page struct {
	Title string
	Body  []byte
}

func (p *page) save() error {
	filename := p.Title + extension
	return ioutil.WriteFile(filename, p.Body, 0600) //WriteFile is a standard library function that writes a byte slice to a file
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("debug Thiago " + r.URL.Path)

	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1<div>%s</div>", p.Title, p.Body)
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
	p1 := &page{Title: "testPage", Body: []byte("Testing with a simple example")}
	p1.save()

	p2, _ := loadPage("testPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
