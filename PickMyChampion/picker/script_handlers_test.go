package picker

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"

	"google.golang.org/appengine/aetest"
)

func TestScriptHandler(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fail()
	}
	defer inst.Close()

	var jsonStr = []byte(`{"dir_path":"/static/img/champion-portraits/"}`)

	req, err := http.NewRequest("GET", "/scripts/go/listDir", bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		t.Errorf("Error making a request to /scripts/go/listDir/: %s", err.Error())
		t.Fail()
	}
	rec := httptest.NewRecorder()
	r := makeRouter("../")
	r.ServeHTTP(rec, req)

	if 200 != rec.Code {
		t.Errorf("Received error code %s.", rec.Code)
		t.Errorf("Body: %s", rec.Body)
		t.Fail()
	}
}
