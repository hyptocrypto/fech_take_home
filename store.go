package main

import (
	"sync"
)

// Since negative points are not in scope, this represents points that have not been calculated/initialized yet.
// Other approach is an int pointer, but using the constant negates the need for ugly memory allocations when using a nil int.
const UninitializedPoints = -1

type ReceiptStore interface {
	Save(r *Receipt) string
	Get(id string) (*Receipt, bool)
}

type Store struct {
	store map[string]*Receipt
	mux   sync.RWMutex
}

// Save will either return the existing record, or create a new one.
// Additionally we kick off the CalculateReceiptPoints in the background.
// The idea here is that this could be a costly operation, and we should not wait on it.
func (s *Store) Save(r *Receipt) string {
	s.mux.Lock()
	defer s.mux.Unlock()
	id := GenerateUUIDForReceipt(r)
	rec, ok := s.store[id]
	if rec != nil && ok {
		// If for some reason the Receipt exists, but points have not been calculated, calculate the points.
		if rec.Points == UninitializedPoints {
			go CalculateReceiptPoints(rec)
		}
		return id
	}
	s.store[id] = r
	go CalculateReceiptPoints(r)
	return id
}

func (s *Store) Get(id string) (*Receipt, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	rec, exists := s.store[id]
	return rec, exists
}

func NewReceiptStore() *Store {
	return &Store{store: make(map[string]*Receipt), mux: sync.RWMutex{}}
}
