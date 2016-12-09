package picker

import (
	"net/http"
	"picker/scripts/golang"
)

func init() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("templates"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/scripts/js/", http.StripPrefix("/scripts/js/", http.FileServer(http.Dir("picker/scripts/js"))))
	http.HandleFunc("/scripts/go/", goScriptHandler)
}

func goScriptHandler(w http.ResponseWriter, r *http.Request) {
	golang.ScriptHandler(w,r)
}
