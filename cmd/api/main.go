package main

import (
	"github.com/codegangsta/negroni"

	"github.com/ohninar/api-sliding-window/api"
)

func main() {
	router := api.Router()
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")
}
