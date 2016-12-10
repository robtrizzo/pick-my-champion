package picker

import (
	"net/http"
	"strings"
	"path/filepath"
	"io/ioutil"
	"fmt"

	"github.com/gorilla/mux"
	"path"
)


func jsScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(".", r.RequestURI))
}


func goScriptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	base := params["func"]
	//map the base name to a function
	//expand this list as functions are added
	m := map[string]func(http.ResponseWriter, *http.Request){
		"listDir": listDir,
	}

	if f, ok := m[base]; ok {
		f(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func listDir(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir_path")
	//fmt.Fprint(w, r.Form)

	if dir == "" {
		http.Error(w, "Could not find dir_path in the form body.", http.StatusNotFound)
		return
	}
	if !strings.HasPrefix(dir, "/") {
		http.Error(w, "The passed dir_path does not start with \"/\".", http.StatusNotImplemented)
		return
	}
	topLevel := strings.Split(dir, "/")[1]
	if topLevel != "static" {
		http.Error(w, "Content can only be listed if a child of the static directory (/static.  Don't forget the " +
				"leading \"/\"). Passed dir: " + topLevel,
			http.StatusNotImplemented)
		return
	}
	final_dir := strings.TrimPrefix(filepath.FromSlash(dir), string(filepath.Separator))
	abs, _ := filepath.Abs(final_dir)

	if !strings.HasPrefix(abs, filepath.Join(wd, "static")) {
		http.Error(w, "Attempted to read outside the static directory.", http.StatusInternalServerError)
		return
	}
	files, err := ioutil.ReadDir(final_dir)
	if err != nil {
		http.Error(w, "Error reading dir " + final_dir + ": " + err.Error(), http.StatusInternalServerError)
	}
	for i, fn := range files {
		if i != 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprint(w, fn.Name())
	}
}
