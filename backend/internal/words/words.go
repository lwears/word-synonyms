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

func (s *WordService) AddWord(word string) (int64, error) {
	result, err := s.DB.ExecContext(
		context.Background(),
		InsertWord, word,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *WordService) FindWord(word string) (WordDbRow, error) {
	var wordRow WordDbRow
	row := s.DB.QueryRowContext(
		context.Background(),
		FindWord, word,
	)
	if err := row.Scan(&wordRow.ID, &wordRow.Word); err != nil {
		return wordRow, err
	}

	return wordRow, nil
}
