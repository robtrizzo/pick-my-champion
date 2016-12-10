package picker

import (
	"net/http"
	"strings"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"path"
	"encoding/json"
	"bytes"

	"github.com/gorilla/mux"
)


func jsScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(appPath, r.URL.Path))
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
	if dir == "" {
		// If we couldn't find the value from the form, look in the straight up body
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		newStr := buf.String()
		var varmap map[string]interface{}
		err := json.Unmarshal([]byte(newStr), &varmap)
		if err != nil {
			http.Error(w, "Error decoding the dir_path from the http request body: " + err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		dir = varmap["dir_path"].(string)
	}

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
	relative_dir := strings.TrimPrefix(filepath.FromSlash(dir), string(filepath.Separator))
	abs, _ := filepath.Abs(filepath.Join(appPath, relative_dir))

	if !strings.HasPrefix(abs, filepath.Join(appPath, "static")) {
		http.Error(w, "Attempted to read outside the static directory.", http.StatusInternalServerError)
		return
	}
	files, err := ioutil.ReadDir(abs)
	if err != nil {
		http.Error(w, "Error reading dir " + abs + ": " + err.Error(), http.StatusInternalServerError)
	}
	for i, fn := range files {
		if i != 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprint(w, fn.Name())
	}
}
