package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/ohninar/api-sliding-window/api"
)

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func TestUploadImageWithImageInvalid(t *testing.T) {
	data := url.Values{}
	data.Set("uploadfile", "teste")

	r, err := http.NewRequest("POST", "/upload", bytes.NewBufferString(data.Encode()))
	ok(t, err)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, r)

	equals(t, http.StatusBadRequest, w.Code)
	equals(t, true, isJSON(w.Body.String()))
	equals(t, "{\"message\":\"request Content-Type isn't multipart/form-data\",\"errors\":null}\n", w.Body.String())
}

func TestUploadImageWithImageValid(t *testing.T) {
	path := "../fixtures/img1.png"

	file, err := os.Open(path)
	ok(t, err)

	fi, _ := file.Stat()

	fileContents, err := ioutil.ReadAll(file)
	ok(t, err)
	file.Close()

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	part, err := bodyWriter.CreateFormFile("uploadfile", fi.Name())
	ok(t, err)

	_, err = part.Write(fileContents)
	ok(t, err)
	bodyWriter.Close()

	r, err := http.NewRequest("POST", "/upload", bodyBuf)
	ok(t, err)
	r.Header.Add("Content-Type", bodyWriter.FormDataContentType())

	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, r)

	equals(t, http.StatusOK, w.Code)
	equals(t, true, isJSON(w.Body.String()))
	equals(t, "{\"message\":\"File processed with success. File name: img1.png (0,0)-(32,39) total sliding=1\"}\n", w.Body.String())
}
