package picker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"path"

	"google.golang.org/appengine/aetest"
)

func TestMain(m *testing.M) {
	startAppEngine()
	code := m.Run()
	os.Exit(code)
}

func startAppEngine() {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		panic(err)
	}
	defer inst.Close()
}


func TestGoChampionDropped(t *testing.T) {
	makeChampionDroppedCall("Aatrox", "team1pick1", http.StatusOK, "", t)
	makeChampionDroppedCall("Aatrox", "", http.StatusInternalServerError, "", t)
	makeChampionDroppedCall("", "team1pick1", http.StatusInternalServerError, "", t)

	jsonStr := "{\"champion_name\":\"Aatrox\"}"
	makeGoCall(jsonStr, "championDropped", http.StatusInternalServerError, "", t)

	jsonStr = "{\"droppable_name\":\"team1pick1\"}"
	makeGoCall(jsonStr, "championDropped", http.StatusInternalServerError, "", t)

	// No comma
	jsonStr = "{\"champion_name\":\"Aatrox\" \"droppable_name\":\"team1pick1\"}"
	makeGoCall(jsonStr, "championDropped", http.StatusInternalServerError, "", t)

	// A comma but no second field
	jsonStr = "{\"champion_name\":\"Aatrox\",}"
	makeGoCall(jsonStr, "championDropped", http.StatusInternalServerError, "", t)

	// No first field
	jsonStr = "{,\"droppable_name\":\"team1pick1\"}"
	makeGoCall(jsonStr, "championDropped", http.StatusInternalServerError, "", t)
}

func makeChampionDroppedCall(championName, droppableName string, expectedCode int, expectedString string, t *testing.T) {
	jsonStr := "{\"champion_name\":\"" + championName + "\", \"droppable_name\":\"" + droppableName + "\"}"
	makeGoCall(jsonStr, "championDropped", expectedCode, expectedString, t)
}


func TestBadGoScript(t *testing.T) {
	makeGETRequest("/scripts/go/someBadFunc", http.StatusNotFound, "", t)
}

func TestFileServingHandler(t *testing.T) {
	makeGETRequest("/scripts/js/champion_areas.js", http.StatusOK, "loadPortraits", t)
	makeGETRequest("/scripts/js/champion_areas.js/anotherdir", http.StatusNotFound, "", t)
	makeGETRequest("/scripts/js/champion_areas", http.StatusNotFound, "", t)
	makeGETRequest("/scripts/js/not_found", http.StatusNotFound, "", t)
	makeGETRequest("/scripts/js/", http.StatusOK, "champion_areas", t)
	makeGETRequest("/scripts/js", http.StatusNotFound, "", t)
	makeGETRequest("/", http.StatusOK, "Pick My Champion", t)
	makeGETRequest("/picker", http.StatusNotFound, "", t)
	makeGETRequest("/static/", http.StatusOK, "css", t)
	makeGETRequest("/static/css/", http.StatusOK, "picker", t)
	makeGETRequest("/static/css/notfound", http.StatusNotFound, "", t)
	makeGETRequest("/static/css/notfound.css", http.StatusNotFound, "", t)
	makeGETRequest("/static/css/picker.css", http.StatusOK, "championPortrait", t)
	makeGETRequest("/static/css/picker.css/anotherdir", http.StatusNotFound, "", t)
	makeGETRequest("/static/img/", http.StatusOK, "champion-portraits", t)
	makeGETRequest("/static/img/champion-portraits/", http.StatusOK, ".png", t)
	makeGETRequest("/static/img/champion-portraits/Aatrox.png", http.StatusOK, "", t)
}

func makeGETRequest(path string, expectedCode int, expectedString string, t *testing.T) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Errorf("Error making a request to %s: %s", path, err.Error())
		t.Fail()
	}
	rec := httptest.NewRecorder()
	r := makeRouter("../")
	r.ServeHTTP(rec, req)

	checkRecorder(expectedCode, expectedString, rec, t)
}

func TestGoListDirHandler(t *testing.T) {
	makeListDirCall("/static/img/champion-portraits/", http.StatusOK, "Aatrox", t)
	makeListDirCall("/static/img/", http.StatusOK, "champion-portraits", t)
	makeListDirCall("static/img/", http.StatusNotFound, "", t)
	makeListDirCall("/", http.StatusNotFound, "", t)
	makeListDirCall("/scripts", http.StatusNotFound, "", t)
	makeListDirCall("/path/to/unknown/dir", http.StatusNotFound, "", t)
	makeListDirCall("/static/img/badpath", http.StatusNotFound, "", t)
	makeListDirCall("/static/img/low/bad/subpath", http.StatusNotFound, "", t)
	makeListDirCall("/static/../../", http.StatusNotFound, "", t)
	makeListDirCall("/app.yaml", http.StatusNotFound, "", t)
	makeListDirCall("", http.StatusInternalServerError, "", t)

	// Test malformed JSON (missing ending })
	makeGoCall("{\"dir_path\":\"\"", "listDir", http.StatusInternalServerError, "", t)
	// Test malformed JSON (empty JSON)
	makeGoCall("{}", "listDir", http.StatusInternalServerError, "", t)
}

func makeListDirCall(path string, expectedCode int, expectedString string, t *testing.T) {
	jsonStr := "{\"dir_path\":\"" + path + "\"}"
	makeGoCall(jsonStr, "listDir", expectedCode, expectedString, t)
}


func makeGoCall(json, funcName string, expectedCode int, expectedString string, t *testing.T) {
	var jsonBytes = []byte(json)

	reqPath := path.Join("/scripts/go", funcName)

	req, err := http.NewRequest("POST", reqPath, bytes.NewBuffer(jsonBytes))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		t.Errorf("Error making a request to /%s: %s", reqPath, err.Error())
		t.Fail()
	}
	rec := httptest.NewRecorder()
	r := makeRouter("../")
	r.ServeHTTP(rec, req)

	checkRecorder(expectedCode, expectedString, rec, t)
}

func checkReturnCode(expectedCode int, recorder *httptest.ResponseRecorder, t *testing.T) {
	if expectedCode != -1 && expectedCode != recorder.Code {
		t.Errorf("Received error code %s.", recorder.Code)
		t.Errorf("Received body: %s", recorder.Body)
		t.Fail()
	}
}

func checkBody(expectedString string, recorder *httptest.ResponseRecorder, t *testing.T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(recorder.Body)
	if expectedString != "" && !strings.Contains(buf.String(), expectedString) {
		t.Errorf("The body of the returned function did not contain an expected string (%s)", expectedString)
		t.Errorf("Received body: %s", recorder.Body)
		t.Fail()
	}
}

func checkRecorder(expectedCode int, expectedString string, recorder *httptest.ResponseRecorder, t *testing.T) {
	checkReturnCode(expectedCode, recorder, t)
	checkBody(expectedString, recorder, t)
}
