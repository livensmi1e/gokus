package room

import (
	"context"
	"sync"
)

type Registry struct {
	mu        sync.Mutex
	rooms     map[string]*Room
	closed    bool
	closeOnce sync.Once
}

func NewRegistry() *Registry {
	return &Registry{rooms: make(map[string]*Room)}
}

func (r *Registry) Join(ctx context.Context, id string) (*Client, error) {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil, ErrClosed
	}
	room, ok := r.rooms[id]
	if !ok {
		// context control lifetime of Room, lives until Room.Close or Client.Leave
		room = New(context.Background())
		r.rooms[id] = room
	}
	// if use defer here, a slow join of previous client will block next client join
	r.mu.Unlock()
	// pass caller's context to join operation, ctx control lifetime of join operation
	return room.Join(ctx)
}

func (r *Registry) Close() {
	r.closeOnce.Do(func() {
		r.mu.Lock()
		r.closed = true
		rooms := make([]*Room, 0, len(r.rooms))
		for _, rm := range r.rooms {
			rooms = append(rooms, rm)
		}
		clear(r.rooms)
		r.mu.Unlock()
		for _, rm := range rooms {
			rm.Close()
		}
	})
}
