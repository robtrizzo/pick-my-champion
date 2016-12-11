package picker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"errors"
	"strings"

	"github.com/gorilla/mux"
)

func goScriptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	base := params["func"]
	//map the base name to a function
	//expand this list as functions are added
	m := map[string]func(http.ResponseWriter, *http.Request){
		"listDir": listDir,
		"championDropped": championDropped,
	}

	if f, ok := m[base]; ok {
		f(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func championDropped(w http.ResponseWriter, r *http.Request) {
	champion_name := r.FormValue("champion_name")
	droppable_name := r.FormValue("droppable_name")
	if champion_name == "" || droppable_name == "" {
		varmap, err := getVarMapFromBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		champstr, ok := varmap["champion_name"].(string)
		if !ok || champstr == "" {
			http.Error(w, "Could not find champion_name in the form body.", http.StatusInternalServerError)
			return
		}
		champion_name = champstr
		droppablestr, ok := varmap["droppable_name"].(string)
		if !ok || droppablestr == "" {
			http.Error(w, "Could not find droppable_name in the form body.", http.StatusInternalServerError)
			return
		}
		droppable_name = droppablestr
	}
	fmt.Fprint(w, "Dropped " + champion_name + " into " + droppable_name)
}

func getVarMapFromBody(r *http.Request) (map[string]interface{}, error) {
	// If we couldn't find the value from the form, look in the straight up body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	newStr := buf.String()
	var varmap map[string]interface{}
	err := json.Unmarshal([]byte(newStr), &varmap)
	if err != nil {
		return nil, errors.New("Error decoding the json from the http request body: " + err.Error())
	}
	defer r.Body.Close()
	return varmap, nil
}

func listDir(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir_path")
	if dir == "" {
		varmap, err := getVarMapFromBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dirstr, ok := varmap["dir_path"].(string)
		if !ok {
			http.Error(w, "Could not find dir_path in the form body.", http.StatusInternalServerError)
			return
		}
		dir = dirstr
	}

	if dir == "" {
		http.Error(w, "Could not find dir_path in the form body.", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(dir, "/") {
		http.Error(w, "The passed dir_path does not start with \"/\".", http.StatusNotFound)
		return
	}
	topLevel := strings.Split(dir, "/")[1]
	if topLevel != "static" {
		http.Error(w, "Content can only be listed if a child of the static directory (/static.  Don't forget the " +
				"leading \"/\"). Passed dir: " + topLevel,
			http.StatusNotFound)
		return
	}
	relative_dir := strings.TrimPrefix(filepath.FromSlash(dir), string(filepath.Separator))
	abs, _ := filepath.Abs(filepath.Join(appPath, relative_dir))

	if !strings.HasPrefix(abs, filepath.Join(appPath, "static")) {
		http.Error(w, "Attempted to read outside the static directory.", http.StatusNotFound)
		return
	}
	files, err := ioutil.ReadDir(abs)
	if err != nil {
		http.Error(w, "Error reading dir " + abs + ": " + err.Error(), http.StatusNotFound)
		return
	}
	for i, fn := range files {
		if i != 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprint(w, fn.Name())
	}
}
