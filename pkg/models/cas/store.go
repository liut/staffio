package cas

type TicketStore interface {
	GetTicket(value string) (*Ticket, error)
	DeleteTicket(value string) error
	SaveTicket(t *Ticket) error
}
