package words

const (
	InsertWord    = `INSERT INTO Words (word) VALUES (?);`
	InsertSynonym = `INSERT INTO Synonyms (word_id, synonym_id) VALUES (?,?) `
	FindWord      = `SELECT * FROM Words WHERE word=?;`
	GetAll        = `SELECT word FROM Words`
	GetWordById   = `SELECT * FROM Words WHERE id=?;`
	GetSynonyms   = `	SELECT Synonyms.word_id, Synonyms.synonym_id, Words.word 
						FROM Synonyms 
						LEFT JOIN Words 
						ON Synonyms.synonym_id = Words.id
						WHERE Synonyms.word_id=?;`
	GetWordsForSynonym = `	
					SELECT Synonyms.word_id, Synonyms.synonym_id, Words.word 
					FROM Synonyms 
					LEFT JOIN Words 
					ON Synonyms.word_id = Words.id
					WHERE Synonyms.synonym_id=?;`
)
