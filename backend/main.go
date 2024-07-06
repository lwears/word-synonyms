package main

import (
	"fmt"
	"net/http"

	"github.com/lwears/word-synonyms/internal/database"
	"github.com/lwears/word-synonyms/internal/words"
)

func main() {
	// This should be passed in via env, but for the sake of simplicity ill add it here
	dbPath := "app.db"
	database, err := database.ConnectAndInitDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	wordService := words.NewWordService(database)
	wordManager := words.NewWordsHTTPHandler(*wordService)

	mux := makeMux(wordManager)
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", mux)
}

func makeMux(m *words.WordHTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /word", m.AddWordHandler)
	mux.HandleFunc("GET /words", m.GetAllWordsHandler)
	mux.HandleFunc("POST /synonym/{word}", m.AddSynonymHandler)
	mux.HandleFunc("GET /synonyms/{word}", m.GetSynonymsHandler)
	mux.HandleFunc("GET /words/{synonym}", m.GetWordsForSynonymHandler)
	return mux
}
