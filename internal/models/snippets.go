package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', $3)
		RETURNING id
	`
	var id int
	err := m.DB.QueryRow(stmt, s.Title, s.Content, s.Expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
