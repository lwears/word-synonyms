package main

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/lwears/word-synonyms/internal/database"
	"github.com/lwears/word-synonyms/internal/words"
	"github.com/rs/cors"
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
	handler := cors.Default().Handler(mux)
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", enforceJSONHandler(handler))
}

func makeMux(m *words.WordHTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /word", m.AddWordHandler)
	mux.HandleFunc("GET /words", m.GetAllWordsHandler)
	mux.HandleFunc("GET /words/{synonym}", m.GetWordsForSynonymHandler)

	mux.HandleFunc("POST /synonym/{word}", m.AddSynonymHandler)
	mux.HandleFunc("GET /synonyms/{word}", m.GetSynonymsHandler)

	return mux
}

// https://www.alexedwards.net/blog/making-and-using-middleware
func enforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
