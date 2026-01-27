package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore stores clips persistently (only pinned items)
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite store
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create table if not exists
	if err := createSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func createSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS clips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		type TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		created_at INTEGER DEFAULT (strftime('%s', 'now'))
	);

	CREATE INDEX IF NOT EXISTS idx_timestamp ON clips(timestamp);
	CREATE INDEX IF NOT EXISTS idx_content ON clips(content);
	`
	_, err := db.Exec(query)
	return err
}

// Add stores a clip (assumed to be pinned)
func (s *SQLiteStore) Add(clip Clip) error {
	query := `
	INSERT INTO clips (content, type, timestamp)
	VALUES (?, ?, ?)
	`
	_, err := s.db.Exec(query, clip.Content, clip.Type, clip.Timestamp.Unix())
	return err
}

// Get retrieves a clip by ID
func (s *SQLiteStore) Get(id int64) (Clip, error) {
	query := `
	SELECT id, content, type, timestamp
	FROM clips
	WHERE id = ?
	`
	row := s.db.QueryRow(query, id)

	var clip Clip
	var ts int64
	err := row.Scan(&clip.ID, &clip.Content, &clip.Type, &ts)
	if err != nil {
		return Clip{}, err
	}

	clip.Timestamp = time.Unix(ts, 0)
	clip.Pinned = true // Everything in SQLite is pinned
	return clip, nil
}

// List returns all pinned clips (most recent first)
func (s *SQLiteStore) List() ([]Clip, error) {
	query := `
	SELECT id, content, type, timestamp
	FROM clips
	ORDER BY timestamp DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []Clip
	for rows.Next() {
		var clip Clip
		var ts int64
		err := rows.Scan(&clip.ID, &clip.Content, &clip.Type, &ts)
		if err != nil {
			return nil, err
		}

		clip.Timestamp = time.Unix(ts, 0)
		clip.Pinned = true
		clips = append(clips, clip)
	}

	return clips, nil
}

// Remove deletes a clip
func (s *SQLiteStore) Remove(id int64) error {
	query := `DELETE FROM clips WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// Search clips by content
func (s *SQLiteStore) Search(query string, limit int) ([]Clip, error) {
	sqlQuery := `
	SELECT id, content, type, timestamp
	FROM clips
	WHERE content LIKE ?
	ORDER BY timestamp DESC
	LIMIT ?
	`
	rows, err := s.db.Query(sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []Clip
	for rows.Next() {
		var clip Clip
		var ts int64
		err := rows.Scan(&clip.ID, &clip.Content, &clip.Type, &ts)
		if err != nil {
			return nil, err
		}

		clip.Timestamp = time.Unix(ts, 0)
		clip.Pinned = true
		clips = append(clips, clip)
	}

	return clips, nil
}

// Count returns the number of pinned clips
func (s *SQLiteStore) Count() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM clips").Scan(&count)
	return count, err
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
