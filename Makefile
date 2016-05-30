install-dep-prod:
	go get github.com/codegangsta/negroni
	go get github.com/gorilla/mux

run:
	go run cmd/api/main.go

test:
	go test -cover github.com/ohninar/api-sliding-window/api

test-cover:
	go test -cover github.com/ohninar/api-sliding-window/api -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
