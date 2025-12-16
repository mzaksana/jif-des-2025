package event

import (
	"log"
	"sync"
)

// Event types
const (
	PostCreated = "post.created"
	PostUpdated = "post.updated"
	PostDeleted = "post.deleted"
)

// Post represents a blog post for events
type Post struct {
	ID        string
	Title     string
	Content   string
	Author    string
	Tags      []string
	CreatedAt int64
	UpdatedAt int64
}

// Event represents a domain event
type Event struct {
	Type string
	Data Post
}

// Handler is a function that handles events
type Handler func(Event)

// Bus is a simple in-memory event bus
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// NewBus creates a new event bus
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for an event type
func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	log.Printf("[EVENT BUS] Publishing event: %s for post ID: %s", e.Type, e.Data.ID)

	if handlers, ok := b.handlers[e.Type]; ok {
		for _, handler := range handlers {
			go handler(e)
		}
	}
}
