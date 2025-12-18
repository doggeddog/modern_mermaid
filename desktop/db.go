package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Diagram struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	IsDeleted bool      `json:"isDeleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table
	query := `
    CREATE TABLE IF NOT EXISTS diagrams (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        title TEXT,
        source TEXT,
        is_deleted BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Helper to extract title from first 2 lines
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	var validLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			validLines = append(validLines, line)
			if len(validLines) >= 2 {
				break
			}
		}
	}
	if len(validLines) == 0 {
		return "Untitled Diagram"
	}
	return strings.Join(validLines, " - ")
}

func (a *App) dbInsertDiagram(content, source string) (int64, error) {
	title := extractTitle(content)
	query := `INSERT INTO diagrams (content, title, source, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`
	res, err := a.db.Exec(query, content, title, source)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) dbUpdateDiagram(id int64, content string) error {
	title := extractTitle(content)
	query := `UPDATE diagrams SET content = ?, title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := a.db.Exec(query, content, title, id)
	return err
}

func (a *App) dbGetLatestDiagram() (*Diagram, error) {
	query := `SELECT id, content, title, source, created_at, updated_at FROM diagrams WHERE is_deleted = 0 ORDER BY updated_at DESC LIMIT 1`
	row := a.db.QueryRow(query)

	var d Diagram
	err := row.Scan(&d.ID, &d.Content, &d.Title, &d.Source, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Optional: For future history list
func (a *App) dbListDiagrams(limit int) ([]Diagram, error) {
	query := `SELECT id, content, title, source, created_at, updated_at FROM diagrams WHERE is_deleted = 0 ORDER BY updated_at DESC LIMIT ?`
	rows, err := a.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diagrams []Diagram
	for rows.Next() {
		var d Diagram
		if err := rows.Scan(&d.ID, &d.Content, &d.Title, &d.Source, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		diagrams = append(diagrams, d)
	}
	return diagrams, nil
}
