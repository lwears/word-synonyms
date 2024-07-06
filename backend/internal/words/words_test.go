package words_test

import (
	"database/sql"
	"testing"

	"github.com/lwears/word-synonyms/internal/database"
	"github.com/lwears/word-synonyms/internal/words"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := database.ConnectAndInitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %s", err)
	}
	return db
}

func TestAdd(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	wordService := words.NewWordService(db)

	word := "funny"
	result, err := wordService.AddWord(word)
	if err != nil {
		t.Errorf("Failed to add word: %s", word)
	}
	if result.ID != 1 {
		t.Errorf("expected row id to be %v got %v", 1, result)
	}

	retrievedWord, err := wordService.GetWord(word)
	if err != nil {
		t.Fatalf("Failed to find word: %s", word)
	}

	// Check if the retrieved word matches the added word
	if retrievedWord.Word != word {
		t.Errorf("expected retrieved word to be '%s', got '%s'", word, retrievedWord.Word)
	}

	if retrievedWord.ID != result.ID {
		t.Errorf("expected retrieved word ID to be '%v', got '%v'", result.ID, retrievedWord.ID)
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
		t.Errorf("Failed to add word: %s", word)
	}
	retrievedWord, err := wordService.GetWord(word)
	if err != nil {
		t.Fatalf("Error finding word: %s", word)
	}
	if retrievedWord == nil {
		t.Errorf("Word does not exist: %s", word)
	}
}

func TestAddAndGetSynonym(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	wordService := words.NewWordService(db)

	word := "dark"
	synonym := "shadowy"

	wordDbRow, err := wordService.GetOrAddWord(word)
	if err != nil {
		t.Errorf("Failed get word: %s", word)
	}

	synonymWordDbRow, err := wordService.GetOrAddWord(synonym)
	if err != nil {
		t.Errorf("Failed get word: %s", word)
	}

	result, err := wordService.AddSynonym(wordDbRow.ID, synonymWordDbRow.ID)
	if err != nil {
		t.Errorf("Failed to add synonym: %s", err)
	}
	if result != 1 {
		t.Errorf("expected row id to be %v got %v", 1, result)
	}
	synonyms, err := wordService.GetSynonyms(wordDbRow)
	if err != nil {
		t.Errorf("Failed to get synonyms: %s", err.Error())
	}
	if synonyms.Word != wordDbRow.Word {
		t.Errorf("Incorrect word returned %s", synonyms.Word)
	}
	if synonyms.Synonyms[0] != synonym {
		t.Errorf("Incorrect synonym returned %s", synonyms.Synonyms[0])
	}
}

func TestGetWordsForSynonyms(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	wordService := words.NewWordService(db)

	word := "dark"
	synonym := "shadowy"

	wordDbRow, err := wordService.GetOrAddWord(word)
	if err != nil {
		t.Errorf("Failed get word: %s", word)
	}

	synonymWordDbRow, err := wordService.GetOrAddWord(synonym)
	if err != nil {
		t.Errorf("Failed get or add word: %s", word)
	}

	result, err := wordService.AddSynonym(wordDbRow.ID, synonymWordDbRow.ID)
	if err != nil {
		t.Errorf("Failed to add synonym: %s", word)
	}
	if result != 1 {
		t.Errorf("expected row id to be %v got %v", 1, result)
	}

	wordsForSynonym, err := wordService.GetWordsForSynonym(synonymWordDbRow)
	if err != nil {
		t.Errorf("Failed to get synonyms: %s", word)
	}
	if wordsForSynonym.Synonym != synonymWordDbRow.Word {
		t.Errorf("Incorrect synonym returned %s", wordsForSynonym.Synonym)
	}
	if wordsForSynonym.Words[0] != wordDbRow.Word {
		t.Errorf("Incorrect synonym returned %s", wordsForSynonym.Words[0])
	}
}
