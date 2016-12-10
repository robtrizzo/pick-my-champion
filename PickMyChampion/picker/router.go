package picker

import (
	"net/http"

	"github.com/gorilla/mux"
	"path/filepath"
	"os"
)
var wd string

func init() {
	http.Handle("/", makeRouter("."))
}

func makeRouter(parentPath string) *mux.Router {
	r := mux.NewRouter()
	wd, _ = os.Getwd()
	wd = filepath.Join(wd, parentPath)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(wd, "static")))))
	r.HandleFunc("/scripts/{lang}/{func}", scriptHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(filepath.Join(wd, "templates")))))
	return r
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lang := params["lang"]
	if lang == "go" {
		goScriptHandler(w, r)
	} else if lang == "js"{
		jsScriptHandler(w, r)
	}
}
