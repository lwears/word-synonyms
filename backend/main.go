package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := makeMux()
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", mux)
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleHello)
	return mux
}

func handleHello(writer http.ResponseWriter,
	request *http.Request,
) {
	fmt.Fprintf(writer, "Hello world!\n")
}
