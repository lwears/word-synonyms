package words

const (
	InsertWord    = `INSERT INTO Words (word) VALUES (?);`
	InsertSynonym = `INSERT INTO Synonyms (synonym_id, word_id) VALUES (?,?) `
	FindWord      = `SELECT * FROM Words WHERE word=?;`
	GetWordById   = `SELECT * FROM Words WHERE id=?;`
	GetSynonyms   = `SELECT * FROM Synonyms 
						LEFT JOIN Words ON Synonyms.synonym_id = Words.id
						WHERE Synonyms.word_id = ?;`
)
