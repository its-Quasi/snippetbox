package models

import (
	"database/sql"
	"errors"
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

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > NOW() AND id = $1
	`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `
		SELECT id, title, content, created, expires
		FROM snippets
		ORDER BY id desc
		LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
