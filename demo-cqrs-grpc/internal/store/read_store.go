package store

import (
	"log"
	"strings"
	"sync"

	"demo-cqrs-grpc/internal/event"
)

// Post represents a denormalized post for fast reads
type Post struct {
	ID        string
	Title     string
	Content   string
	Author    string
	Tags      []string
	CreatedAt int64
	UpdatedAt int64
}

// ReadStore handles read operations from in-memory store
// In production, this would be Redis or Elasticsearch
type ReadStore struct {
	mu    sync.RWMutex
	posts map[string]*Post
}

// NewReadStore creates a new read store
func NewReadStore(bus *event.Bus) *ReadStore {
	rs := &ReadStore{
		posts: make(map[string]*Post),
	}

	// Subscribe to events to keep read store in sync
	bus.Subscribe(event.PostCreated, rs.handlePostCreated)
	bus.Subscribe(event.PostUpdated, rs.handlePostUpdated)
	bus.Subscribe(event.PostDeleted, rs.handlePostDeleted)

	log.Println("[READ STORE] In-memory store initialized and subscribed to events")
	return rs
}

func (s *ReadStore) handlePostCreated(e event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.posts[e.Data.ID] = &Post{
		ID:        e.Data.ID,
		Title:     e.Data.Title,
		Content:   e.Data.Content,
		Author:    e.Data.Author,
		Tags:      e.Data.Tags,
		CreatedAt: e.Data.CreatedAt,
		UpdatedAt: e.Data.UpdatedAt,
	}
	log.Printf("[READ STORE] Synced new post to memory: %s", e.Data.ID)
}

func (s *ReadStore) handlePostUpdated(e event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if post, exists := s.posts[e.Data.ID]; exists {
		post.Title = e.Data.Title
		post.Content = e.Data.Content
		post.Tags = e.Data.Tags
		post.UpdatedAt = e.Data.UpdatedAt
		log.Printf("[READ STORE] Synced updated post to memory: %s", e.Data.ID)
	}
}

func (s *ReadStore) handlePostDeleted(e event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.posts, e.Data.ID)
	log.Printf("[READ STORE] Removed post from memory: %s", e.Data.ID)
}

// GetPost retrieves a post by ID
func (s *ReadStore) GetPost(id string) (*Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if exists {
		log.Printf("[READ STORE] Retrieved post from memory: %s", id)
	}
	return post, exists
}

// ListPosts returns paginated posts
func (s *ReadStore) ListPosts(limit, offset int) ([]*Post, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice
	all := make([]*Post, 0, len(s.posts))
	for _, post := range s.posts {
		all = append(all, post)
	}

	total := len(all)

	// Apply pagination
	if offset >= len(all) {
		return []*Post{}, total
	}

	end := offset + limit
	if end > len(all) {
		end = len(all)
	}

	log.Printf("[READ STORE] Listed %d posts from memory (offset: %d, limit: %d)", end-offset, offset, limit)
	return all[offset:end], total
}

// SearchPosts searches posts by query and tags
func (s *ReadStore) SearchPosts(query string, tags []string) []*Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*Post
	query = strings.ToLower(query)

	for _, post := range s.posts {
		// Check query match in title or content
		queryMatch := query == "" ||
			strings.Contains(strings.ToLower(post.Title), query) ||
			strings.Contains(strings.ToLower(post.Content), query)

		// Check tag match
		tagMatch := len(tags) == 0 || s.hasAnyTag(post.Tags, tags)

		if queryMatch && tagMatch {
			results = append(results, post)
		}
	}

	log.Printf("[READ STORE] Search found %d posts (query: %q, tags: %v)", len(results), query, tags)
	return results
}

func (s *ReadStore) hasAnyTag(postTags, searchTags []string) bool {
	for _, st := range searchTags {
		for _, pt := range postTags {
			if strings.EqualFold(st, pt) {
				return true
			}
		}
	}
	return false
}
