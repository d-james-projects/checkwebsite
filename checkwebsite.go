package main

import (
	"net/http"
)

func main() {

	//Handling the /red
	http.HandleFunc("/red", redHandler)

	//Handling the /green
	aHandler := greenHandler{}
	http.Handle("/green", aHandler)

	http.ListenAndServe(":8789", nil)
}

func redHandler(res http.ResponseWriter, req *http.Request) {
	data := []byte("You asked for RED\n")
	res.WriteHeader(200)
	res.Write(data)
}

type greenHandler struct{}

func (h greenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data := []byte("You asked for GREEN\n")
	res.WriteHeader(200)
	res.Write(data)
}

func checkWebsite() (b bool, err error) {
	b = true
	return
}
