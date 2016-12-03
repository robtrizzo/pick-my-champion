package picker

import (
	//"fmt"
	"net/http"
	"html/template"
)

func init() {
	http.HandleFunc("/", handler)
	//http.Handle("/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
}

type Page struct {
	Title string
	Body  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Hello, world!")
	t, err := template.ParseFiles("resources/index.html")
	if t == nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := &Page{Title: "Hello, world", Body: "This is the body of hello, world"}
	t.Execute(w, p)
}
