package cas

import (
	"sync"
	"time"
)

// MemoryStore is memory based Storage
type MemoryStore struct {
	tickets *sync.Map
}

// NewMemoryStore returns new instance of MemoryStore
func NewInMemory() *MemoryStore {
	return &MemoryStore{
		tickets: &sync.Map{},
	}
}

// DoesTicketExist checks if given ticket exists
func (s *MemoryStore) GetTicket(ticket string) *Ticket {
	if v, ok := s.tickets.Load(ticket); ok {
		// check if ticket should be deleted
		t := v.(*Ticket)
		if t.IsOld() {
			s.DeleteTicket(ticket)
			return nil
		}
		return t
	}

	return nil
}

// SaveTicket stores the Ticket
func (s *MemoryStore) SaveTicket(ticket *Ticket) {
	ticket.CreatedAt = time.Now()
	s.tickets.Store(ticket.Value, ticket)
}

// DeleteTicket deletes given ticket
func (s *MemoryStore) DeleteTicket(ticket string) {
	s.tickets.Delete(ticket)
}
