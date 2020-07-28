package main

import (
	"fmt"
	"io/ioutil"
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
}
