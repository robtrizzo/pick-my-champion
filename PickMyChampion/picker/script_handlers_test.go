package picker

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"google.golang.org/appengine/aetest"
	"fmt"
	"os"
)

func TestScriptHandler(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fail()
	}
	defer inst.Close()

	cwd, _:= os.Getwd()
	fmt.Printf("cwd: %s", cwd)
	//form := url.Values{}
	//form.Add("dir_path", "/static/img/champion-portraits/")

	//b := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequest("POST", "/scripts/go/listDir/", nil)
	if err != nil {
		t.Errorf("Error making a request to /scripts/go/listDir/: %s", err.Error())
		t.Fail()
	}
	rec := httptest.NewRecorder()
	r := makeRouter("../")
	r.ServeHTTP(rec, req)
	//fmt.Print(rec.Body)

	if 200 != rec.Code {
		t.Errorf("Received error code %s.", rec.Code)
		t.Errorf("Body: %s", rec.Body)
		t.Fail()
	}
}
