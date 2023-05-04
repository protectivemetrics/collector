package impression

import "database/sql"

type Impression struct {
	Url string
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(impression Impression) error {
	_, err := r.db.Exec("INSERT INTO impression(url, created) values(?, datetime('now'))", impression.Url)
	return err
}

func (r *SQLiteRepository) Migrate() error {
	// Create the table
	query := []string{
		"CREATE TABLE IF NOT EXISTS impression (url TEXT, created TEXT)",
		"CREATE INDEX IF NOT EXISTS impressionCreated ON impression (created)",
		"CREATE INDEX IF NOT EXISTS impressionUrlCreated ON impression (url, created)",
	}
	// Loop through the queries running them
	for _, element := range query {
		_, err := r.db.Exec(element)
		if err != nil {
			return err
		}
	}
	return nil
}
