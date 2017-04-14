package backends

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lcgc/platform/staffio/models/cas"
)

func TestTicket(t *testing.T) {
	service := "http://localhost:3001"
	uid := "test"
	st := cas.NewTicket("ST", service, uid, true)
	err := SaveTicket(st)
	if err != nil {
		t.Logf("save %s ERR %s", st.Value, err)
	}
	assert.Nil(t, err)
	ticket, err := GetTicket(st.Value)
	assert.Nil(t, err)
	assert.NotEmpty(t, ticket.Service, ticket.Uid)
	assert.NotZero(t, ticket.Id)
	assert.Nil(t, DeleteTicket(st.Value))
}
