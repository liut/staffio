package backends

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liut/staffio/pkg/models/cas"
)

func TestTicket(t *testing.T) {
	service := "http://localhost:3001"
	uid := "test"
	st := cas.NewTicket("ST", service, uid, true)
	err := svc.SaveTicket(st)
	if err != nil {
		t.Logf("save %s ERR %s", st.Value, err)
	}
	assert.Nil(t, err)
	ticket, err := svc.GetTicket(st.Value)
	assert.Nil(t, err)
	assert.NotEmpty(t, ticket.Service, ticket.UID)
	assert.NotZero(t, ticket.Id)
	assert.Nil(t, svc.DeleteTicket(st.Value))
}
