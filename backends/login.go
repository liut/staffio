package backends

import (
	"fmt"
	"tuluu.com/liut/staffio/backends/ldap"
	"tuluu.com/liut/staffio/models"
)

func Login(username, password string) (*models.Staff, error) {
	r, staff := ldap.Login(username, password)
	if r {
		return staff, nil
	}
	return nil, fmt.Errorf("Login failed %s", username)
}
