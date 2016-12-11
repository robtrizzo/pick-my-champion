package picker

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var appPath string

func init() {
	http.Handle("/", makeRouter("."))
}

func makeRouter(parentPath string) *mux.Router {
	r := mux.NewRouter()
	appPath, _ = os.Getwd()
	appPath = filepath.Join(appPath, parentPath)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(appPath, "static")))))
	r.PathPrefix("/scripts/js/").Handler(http.StripPrefix("/scripts/js/", http.FileServer(http.Dir(filepath.Join(appPath, "scripts/js")))))
	//r.HandleFunc("/scripts/js/{func}", jsScriptHandler)
	r.HandleFunc("/scripts/go/{func}", goScriptHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(filepath.Join(appPath, "templates")))))
	return r
}
