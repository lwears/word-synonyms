package words

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type WordService struct {
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

func (s *WordService) AddWord(word string) (*WordDbRow, error) {
	result, err := s.DB.ExecContext(
		context.Background(),
		InsertWord, word,
	)
	if err != nil {
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

// Want to be able to return nil in event of no rows, instead of an error, for this i need a pointer
func (s *WordService) GetWord(word string) (*WordDbRow, error) {
	var wordRow WordDbRow
	row := s.DB.QueryRowContext(
		context.Background(),
		FindWord, word,
	)
	if err := row.Scan(&wordRow.ID, &wordRow.Word); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil when no rows are found
		}
		return nil, err // Return other errors
	}

	return &wordRow, nil
}

func (s *WordService) GetWordById(id int64) (*WordDbRow, error) {
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
