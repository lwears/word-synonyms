package words

const (
	InsertWord = `INSERT INTO Words (word) VALUES (?);`
	FindWord   = `SELECT * FROM Words WHERE word=?;`
)
