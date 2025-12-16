package store

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"demo-cqrs-grpc/internal/event"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// WriteStore handles write operations to SQLite
type WriteStore struct {
	db  *sql.DB
	bus *event.Bus
}

// NewWriteStore creates a new write store with SQLite
func NewWriteStore(bus *event.Bus) (*WriteStore, error) {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		return nil, err
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			tags TEXT,
			created_at INTEGER,
			updated_at INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	log.Println("[WRITE STORE] SQLite database initialized")
	return &WriteStore{db: db, bus: bus}, nil
}

// CreatePost creates a new post in SQLite
func (s *WriteStore) CreatePost(title, content, author string, tags []string) (string, error) {
	id := uuid.New().String()
	now := time.Now().Unix()
	tagsStr := strings.Join(tags, ",")

	_, err := s.db.Exec(
		"INSERT INTO posts (id, title, content, author, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, title, content, author, tagsStr, now, now,
	)
	if err != nil {
		return "", err
	}

	log.Printf("[WRITE STORE] Created post in SQLite: %s", id)

	// Publish event to sync with read store
	s.bus.Publish(event.Event{
		Type: event.PostCreated,
		Data: event.Post{
			ID:        id,
			Title:     title,
			Content:   content,
			Author:    author,
			Tags:      tags,
			CreatedAt: now,
			UpdatedAt: now,
		},
	})

	return id, nil
}

// UpdatePost updates an existing post
func (s *WriteStore) UpdatePost(id, title, content string, tags []string) error {
	now := time.Now().Unix()
	tagsStr := strings.Join(tags, ",")

	result, err := s.db.Exec(
		"UPDATE posts SET title = ?, content = ?, tags = ?, updated_at = ? WHERE id = ?",
		title, content, tagsStr, now, id,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Printf("[WRITE STORE] Updated post in SQLite: %s", id)

	// Get full post for event
	var author string
	var createdAt int64
	err = s.db.QueryRow("SELECT author, created_at FROM posts WHERE id = ?", id).Scan(&author, &createdAt)
	if err != nil {
		return err
	}

	s.bus.Publish(event.Event{
		Type: event.PostUpdated,
		Data: event.Post{
			ID:        id,
			Title:     title,
			Content:   content,
			Author:    author,
			Tags:      tags,
			CreatedAt: createdAt,
			UpdatedAt: now,
		},
	})

	return nil
}

// DeletePost deletes a post
func (s *WriteStore) DeletePost(id string) error {
	result, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Printf("[WRITE STORE] Deleted post from SQLite: %s", id)

	s.bus.Publish(event.Event{
		Type: event.PostDeleted,
		Data: event.Post{ID: id},
	})

	return nil
}

// Close closes the database connection
func (s *WriteStore) Close() error {
	return s.db.Close()
}
