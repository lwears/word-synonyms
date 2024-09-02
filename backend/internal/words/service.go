package words

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type WordsService struct {
	DB *sql.DB
}

type WordDbRow struct {
	ID   int64
	Word string
}

type Synonyms struct {
	Word     string
	Synonyms []string
}

type WordsForSynonym struct {
	Synonym string
	Words   []string
}

type GetSynonymsResult struct {
	SynonymID int
	WordID    int
	Word      string
}

func NewWordsService(DB *sql.DB) WordsService {
	return WordsService{
		DB: DB,
	}
}

func (s *WordsService) AddWord(word string) (*WordDbRow, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(InsertWord)
	if err != nil {
		tx.Rollback() // Rollback in case of error
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err := stmt.ExecContext(context.Background(), )
	result, err := s.DB.ExecContext(
		context.Background(),
		InsertWord, word,
	)
	if err != nil {
		fmt.Print(err)
		return &WordDbRow{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &WordDbRow{
		ID:   id,
		Word: word,
	}, nil
}

func (s *WordsService) GetAll() ([]string, error) {
	words := make([]string, 0)
	rows, err := s.DB.QueryContext(
		context.Background(),
		GetAll,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			fmt.Print(err)
			return nil, err
		}
		words = append(words, word)
	}
	return words, err
}

// Want to be able to return nil in event of no rows, instead of an error, for this i need a pointer
// because i think an empty response is better than an error.
func (s *WordsService) GetWord(word string) (*WordDbRow, error) {
	var wordRow WordDbRow
	row := s.DB.QueryRowContext(
		context.Background(),
		FindWord, word,
	)
	if err := row.Scan(&wordRow.ID, &wordRow.Word); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &wordRow, nil
}

func (s *WordsService) GetWordById(id int64) (*WordDbRow, error) {
	var wordRow WordDbRow
	row := s.DB.QueryRowContext(
		context.Background(),
		GetWordById, id,
	)
	if err := row.Scan(&wordRow.ID, &wordRow.Word); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &wordRow, nil
}

func (s *WordsService) GetOrAddWord(word string) (*WordDbRow, error) {
	wordRow, err := s.GetWord(word)
	if err != nil {
		return nil, err
	}
	if wordRow == nil {
		wordRow, err = s.AddWord(word)
		if err != nil {
			return nil, err
		}
	}
	return wordRow, nil
}

func (s *WordsService) AddSynonym(wordId int64, synonymId int64) (int64, error) {
	result, err := s.DB.ExecContext(
		context.Background(),
		InsertSynonym, wordId, synonymId,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// I got stuck here on whether to check for word existence in either server or handler
func (s *WordsService) GetSynonyms(w *WordDbRow) (Synonyms, error) {
	synonyms := Synonyms{
		Word:     w.Word,
		Synonyms: make([]string, 0),
	}

	// fmt.Printf("Executing query with word ID: %d\n", w.ID)
	rows, err := s.DB.QueryContext(context.Background(), GetSynonyms, w.ID)
	if err != nil {
		return synonyms, err
	}

	defer rows.Close()

	for rows.Next() {
		var synonymsResult GetSynonymsResult
		if err := rows.Scan(&synonymsResult.SynonymID, &synonymsResult.WordID, &synonymsResult.Word); err != nil {
			return synonyms, err
		}
		synonyms.Synonyms = append(synonyms.Synonyms, synonymsResult.Word)
	}

	if err := rows.Err(); err != nil {
		return synonyms, fmt.Errorf("error iterating over rows: %v", err)
	}

	return synonyms, nil
}

func (s *WordsService) GetWordsForSynonym(synonym *WordDbRow) (WordsForSynonym, error) {
	wordsForSynonym := WordsForSynonym{
		Synonym: synonym.Word,
		Words:   make([]string, 0),
	}

	rows, err := s.DB.QueryContext(context.Background(), GetWordsForSynonym, synonym.ID)
	if err != nil {
		return wordsForSynonym, err
	}

	defer rows.Close()

	for rows.Next() {
		var synonymsResult GetSynonymsResult
		if err := rows.Scan(&synonymsResult.SynonymID, &synonymsResult.WordID, &synonymsResult.Word); err != nil {
			return wordsForSynonym, err
		}
		wordsForSynonym.Words = append(wordsForSynonym.Words, synonymsResult.Word)
	}

	if err := rows.Err(); err != nil {
		return wordsForSynonym, fmt.Errorf("error iterating over rows: %v", err)
	}

	return wordsForSynonym, nil
}
