package backends

import (
	"errors"
	"log"
	"tuluu.com/liut/staffio/backends/ldap"
	"tuluu.com/liut/staffio/models"
)

var (
	ErrLogin = errors.New("Invalid Username/Password")
)

func Login(username, password string) (*models.Staff, error) {
	r, staff := ldap.Login(username, password)
	if r {
		return staff, nil
	}
	log.Printf("Login failed %s", username)
	return nil, ErrLogin
}
