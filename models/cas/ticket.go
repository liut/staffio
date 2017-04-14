package cas

import (
	"fmt"
	"net/url"
	"time"

	"lcgc/platform/staffio/models/random"
)

const (
	ValueLength    = 128
	MaxValueLength = 139
	MinValueLength = 32
)

type Ticket struct {
	Id        int       `db:"id,pk" json:"id" form:"id"` // seriel in database
	Class     string    `db:"type" json:"type"`          // ticket type: ST, PGT, PT, ...
	Value     string    `db:"value" json:"value"`        // ticket id: (ST-|PGT-|PT-)
	Uid       string    `db:"uid" json:"uid"`            // uid in staff
	Service   string    `db:"service" json:"service"`    // is an URL
	CreatedAt time.Time `db:"created" json:"created"`
	Renew     bool
}

func NewTicket(class string, service string, uid string, renew bool) *Ticket {
	t := Ticket{
		Class:     class,
		Value:     fmt.Sprintf("%s-%s", class, random.GenString(ValueLength)),
		CreatedAt: time.Now(),
		Uid:       uid,
		Service:   service,
		Renew:     renew}
	return &t
}

func (t *Ticket) IsOld() bool {
	return t.CreatedAt.Add(5 * time.Minute).Before(time.Now())
}

func (t *Ticket) Check() *CasError {
	err := ValidateTicket(t.Value)
	if err != nil {
		return err
	}
	if err = ValidateService(t.Service); err != nil {
		return err
	}

	if t.Uid == "" {
		return NewCasError("empty ticket.Uid", ERROR_CODE_INVALID_USERNAME)
	}

	return nil
}

func ValidateService(service string) *CasError {
	if service == "" {
		return NewCasError("empty service url", ERROR_CODE_INVALID_SERVICE)
	}

	_, err := url.Parse(service)
	if err != nil {
		return NewCasError(err.Error(), ERROR_CODE_INVALID_SERVICE)
	}
	return nil
}

func ValidateTicket(ticket string) *CasError {
	err := validateTicketLength(ticket)
	if err != nil {
		return err
	}

	err = validateTicketFormat(ticket)
	if err != nil {
		return err
	}

	return nil
}

func validateTicketLength(ticket string) *CasError {
	if len(ticket) == 0 {
		return NewCasError("Required query parameter `ticket` was not defined.", ERROR_CODE_INVALID_REQUEST)
	}

	if len(ticket) < MinValueLength {
		return NewCasError(fmt.Sprintf(
			"Ticket is not long enough. Minimum length is `%d` but length was `%d`.",
			MinValueLength, len(ticket)), ERROR_CODE_INVALID_TICKET_SPEC)
	}

	if len(ticket) > MaxValueLength {
		return NewCasError(fmt.Sprintf(
			"Ticket is too long. Maximum length is `%d` but length was `%d`.",
			MaxValueLength, len(ticket)), ERROR_CODE_INVALID_TICKET_SPEC)
	}

	return nil
}

func validateTicketFormat(ticket string) *CasError {
	if ticket[0:3] == "ST-" {
		return nil
	} else if ticket[0:4] == "PGT-" {
		return nil
	} else if ticket[0:3] == "PT-" {
		return nil
	}

	return NewCasError("Required ticket prefix is missing. Supported prefixes are: [ST, PGT]",
		ERROR_CODE_INVALID_TICKET_SPEC)
}
