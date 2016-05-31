package api

import "github.com/gorilla/mux"

//Router ...
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/upload", HandleUploadImage).Methods("POST")
	return router
}
