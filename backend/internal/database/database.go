package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const InitialQuery = `
CREATE TABLE IF NOT EXISTS Words (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		word VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS Synonyms (
		word_id INT,
		synonym_id INT,
		PRIMARY KEY (word_id, synonym_id),
		FOREIGN KEY (word_id) REFERENCES Words(word_id),
		FOREIGN KEY (synonym_id) REFERENCES Words(word_id)
);
		`

func ConnectAndInitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		InitialQuery,
	)
	if err != nil {
		return nil, err
	}
	// Check if the database connection is alive
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
