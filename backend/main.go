package main

import (
	"fmt"
	"net/http"

	"github.com/lwears/word-synonyms/internal/database"
)

func main() {
	dbPath := "app.db"
	database, err := database.ConnectAndInitDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	mux := makeMux()
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", mux)
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /word", handleHello)
	// mux.HandleFunc("GET /words", handleHello)
	// mux.HandleFunc("POST /synonym/{word}", handleHello)
	// mux.HandleFunc("GET /synonyms/{word}", handleHello)
	// mux.HandleFunc("GET /words/{synonym}", handleHello)
	return mux
}

func handleHello(writer http.ResponseWriter,
	request *http.Request,
) {
	fmt.Fprintf(writer, "Hello world!\n")
}
