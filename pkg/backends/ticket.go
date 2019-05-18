package backends

import (
	"log"

	"github.com/liut/staffio/pkg/models/cas"
)

func (s *serviceImpl) GetTicket(value string) (*cas.Ticket, error) {
	a := new(cas.Ticket)

	qs := func(db dber) error {
		return db.Get(a, `SELECT id, type, uid, value, service, created FROM cas_ticket WHERE value = $1`, value)
	}
	return a, withDbQuery(qs)
}

func (s *serviceImpl) DeleteTicket(value string) error {
	if value != "" {
		return withTxQuery(func(db dbTxer) error {
			_, err := db.Exec("DELETE from cas_ticket WHERE value = $1", value)
			if err != nil {
				log.Printf("delete ticket %s ERR %s", value, err)
			}
			return err
		})
	}

	return cas.NewCasError("empty ticket value", cas.ERROR_CODE_INVALID_TICKET_SPEC)
}

func (s *serviceImpl) SaveTicket(t *cas.Ticket) error {
	if err := t.Check(); err != nil {
		return err
	}

	return withTxQuery(func(db dbTxer) error {
		_, err := db.Exec("INSERT INTO cas_ticket (type, value, uid, service, created) VALUES($1, $2, $3, $4, $5)",
			t.Class, t.Value, t.UID, t.Service, t.CreatedAt)
		if err != nil {
			log.Printf("save tick %v ERR %s", t, err)
		}
		return err
	})
}
