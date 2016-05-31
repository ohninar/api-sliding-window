package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(status)

	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("respond:", err)
	}
}

//HandleUploadImage ...
func HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	fmt.Println(r.ContentLength)
	err := r.ParseMultipartForm(2000)
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		fmt.Println("ParseMultipartForm", body)
		respond(w, r, http.StatusBadRequest, body)
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		fmt.Println("FormFile", file, handler, body)
		respond(w, r, http.StatusBadRequest, body)
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		body := ErrMessage{Message: err.Error(), Errors: nil}
		respond(w, r, http.StatusBadRequest, body)
		return
	}

	bounds := img.Bounds()

	data := "File processed with success. File name: " + handler.Filename + " " + bounds.String()

	body := SuccessMessage{Message: data}
	respond(w, r, http.StatusOK, body)
}
