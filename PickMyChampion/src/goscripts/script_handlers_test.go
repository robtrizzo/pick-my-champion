package goscripts

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"google.golang.org/appengine/aetest"
)

func TestScriptHandler(t *testing.T) {
	c, done, err := aetest.NewContext()
	if err != nil {
		t.Fail()
	}
	defer done()
	r, _ := http.NewRequest("GET", "/scripts/go/", nil)
	w := httptest.NewRecorder()
	ScriptHandler(w, r)
	c.Value()
	if 200 != w.Code {
		t.Fail()
	}
}
