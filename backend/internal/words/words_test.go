package words_test

import (
	"database/sql"
	"testing"

	"api/internal/database"
	"api/internal/words"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := database.ConnectAndInitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	return db
}

func TestAdd(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	wordService := &words.WordService{
		DB: db,
	}

	word := "funny"
	result, err := wordService.AddWord(word)
	if err != nil {
		t.Errorf("Failed to add word: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected row id to be %v got %v", 1, result)
	}

	retrievedWord, err := wordService.GetWord(word)
	if err != nil {
		t.Fatalf("Failed to find word: %v", err)
	}

	// Check if the retrieved word matches the added word
	if retrievedWord.Word != word {
		t.Errorf("expected retrieved word to be '%v', got '%v'", word, retrievedWord.Word)
	}

	if retrievedWord.ID != result.ID {
		t.Errorf("expected retrieved word ID to be '%v', got '%v'", result, retrievedWord.ID)
	}
}

func TestGetWord(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	wordService := &words.WordService{
		DB: db,
	}

	word := "funny"
	_, err := wordService.AddWord(word)
	if err != nil {
		t.Errorf("Failed to add word: %v", err)
	}
	retrievedWord, err := wordService.GetWord(word)
	if err != nil {
		t.Fatalf("Error finding word: %v", err)
	}
	if retrievedWord == nil {
		t.Errorf("Word does not exist: %v", err)
	}
}
