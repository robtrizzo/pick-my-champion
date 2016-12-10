package picker

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"strings"

	"google.golang.org/appengine/aetest"
)

func TestGOScriptHandler(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fail()
	}
	defer inst.Close()

	makeListDirCall("/static/img/champion-portraits/", http.StatusOK, "Aatrox", t)
	makeListDirCall("/static/img/", http.StatusOK, "champion-portraits", t)
	makeListDirCall("/", http.StatusNotFound, "", t)
	makeListDirCall("/scripts", http.StatusNotFound, "", t)
	makeListDirCall("/path/to/unknown/dir", http.StatusNotFound, "", t)
	makeListDirCall("/static/img/badpath", http.StatusNotFound, "", t)
	makeListDirCall("/static/img/low/bad/subpath", http.StatusNotFound, "", t)
	makeListDirCall("/static/../../", http.StatusNotFound, "", t)
	makeListDirCall("/app.yaml", http.StatusNotFound, "", t)
}

func makeListDirCall(path string, expectedCode int, expectedString string, t *testing.T) {
	var jsonStr = []byte("{\"dir_path\":\"" + path + "\"}")

	req, err := http.NewRequest("POST", "/scripts/go/listDir", bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		t.Errorf("Error making a request to /scripts/go/listDir/: %s", err.Error())
		t.Fail()
	}
	rec := httptest.NewRecorder()
	r := makeRouter("../")
	r.ServeHTTP(rec, req)

	if expectedCode != -1 && expectedCode != rec.Code {
		t.Errorf("Received error code %s.", rec.Code)
		t.Errorf("Received body: %s", rec.Body)
		t.Fail()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rec.Body)
	if expectedString != "" && !strings.Contains(buf.String(), expectedString) {
		t.Errorf("The body of the returned function did not contain an expected string (%s)", expectedString)
		t.Errorf("Received body: %s", rec.Body)
		t.Fail()
	}
}