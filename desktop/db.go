package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	Hashcode  string    `json:"hashcode"`
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
        hashcode TEXT,
        is_deleted BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	// Migration: Check if hashcode column exists
	// If not, add it. (For existing databases from previous step)
	var colName string
	err = db.QueryRow("SELECT name FROM pragma_table_info('diagrams') WHERE name='hashcode'").Scan(&colName)
	if err == sql.ErrNoRows {
		_, err = db.Exec("ALTER TABLE diagrams ADD COLUMN hashcode TEXT")
		if err != nil {
			fmt.Printf("Error adding hashcode column: %v\n", err)
		} else {
			// Backfill hashes
			go backfillHashes(db)
		}
	}

	return db, nil
}

func backfillHashes(db *sql.DB) {
	rows, err := db.Query("SELECT id, content FROM diagrams WHERE hashcode IS NULL")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var content string
		if err := rows.Scan(&id, &content); err == nil {
			hash := calculateHash(content)
			db.Exec("UPDATE diagrams SET hashcode = ? WHERE id = ?", hash, id)
		}
	}
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

func calculateHash(content string) string {
	hasher := md5.New()
	hasher.Write([]byte(content))
	hash := hex.EncodeToString(hasher.Sum(nil))
	if len(hash) > 16 {
		return hash[:16]
	}
	return hash
}

func (a *App) dbInsertDiagram(content, source string) (int64, error) {
	hash := calculateHash(content)

	// Deduplication: Check if exists
	var existingID int64
	err := a.db.QueryRow("SELECT id FROM diagrams WHERE hashcode = ? AND is_deleted = 0", hash).Scan(&existingID)
	if err == nil {
		// Found existing, update updated_at
		a.db.Exec("UPDATE diagrams SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", existingID)
		return existingID, nil
	}

	title := extractTitle(content)
	query := `INSERT INTO diagrams (content, title, source, hashcode, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`
	res, err := a.db.Exec(query, content, title, source, hash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) dbUpdateDiagram(id int64, content string) error {
	title := extractTitle(content)
	hash := calculateHash(content)
	// Also set is_deleted = 0 to resurrect if it was cleared
	query := `UPDATE diagrams SET content = ?, title = ?, hashcode = ?, updated_at = CURRENT_TIMESTAMP, is_deleted = 0 WHERE id = ?`
	_, err := a.db.Exec(query, content, title, hash, id)
	return err
}

func (a *App) dbClearHistory() error {
	_, err := a.db.Exec("UPDATE diagrams SET is_deleted = 1")
	return err
}

func (a *App) dbGetLatestDiagram() (*Diagram, error) {
	query := `SELECT id, content, title, source, hashcode, created_at, updated_at FROM diagrams WHERE is_deleted = 0 ORDER BY updated_at DESC LIMIT 1`
	row := a.db.QueryRow(query)

	var d Diagram
	var hash sql.NullString // Handle potential nulls during migration/init
	err := row.Scan(&d.ID, &d.Content, &d.Title, &d.Source, &hash, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	d.Hashcode = hash.String
	return &d, nil
}

func (a *App) dbGetDiagram(id int64) (*Diagram, error) {
	query := `SELECT id, content, title, source, hashcode, created_at, updated_at FROM diagrams WHERE id = ? AND is_deleted = 0`
	row := a.db.QueryRow(query, id)

	var d Diagram
	var hash sql.NullString
	err := row.Scan(&d.ID, &d.Content, &d.Title, &d.Source, &hash, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	d.Hashcode = hash.String
	return &d, nil
}

// Optional: For future history list
func (a *App) dbListDiagrams(limit int) ([]Diagram, error) {
	query := `SELECT id, content, title, source, hashcode, created_at, updated_at FROM diagrams WHERE is_deleted = 0 ORDER BY updated_at DESC LIMIT ?`
	rows, err := a.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diagrams []Diagram
	for rows.Next() {
		var d Diagram
		var hash sql.NullString
		if err := rows.Scan(&d.ID, &d.Content, &d.Title, &d.Source, &hash, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		d.Hashcode = hash.String
		diagrams = append(diagrams, d)
	}
	return diagrams, nil
}
