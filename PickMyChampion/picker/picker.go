package picker

import (
	"net/http"
)

func init() {
	http.Handle("/picker/", http.StripPrefix("/picker/", http.FileServer(http.Dir("picker/templates"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("picker/scripts"))))
}
