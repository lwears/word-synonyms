package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lwears/word-synonyms/internal/database"
	"github.com/lwears/word-synonyms/internal/words"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	dbPath := os.Getenv("DB_PATH")
	database, err := database.ConnectAndInitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	wordsService := words.NewWordsService(database)
	wordsHandler := words.NewWordsHTTPHandler(wordsService)

	apiHandler := http.StripPrefix("/api", wordsHandler)

	handler := cors.Default().Handler(apiHandler)
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", enforceJSONHandler(handler))
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
