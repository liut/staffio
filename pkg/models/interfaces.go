package models

// StaffStore Storage for Staff
type StaffStore interface {
	// All browse from store, like LDAP
	All() Staffs
	// Get with uid
	Get(uid string) (*Staff, error)
	// GetByDN with dn
	GetByDN(dn string) (*Staff, error)
	// Delete with uid
	Delete(uid string) error
	// Save add or update
	Save(staff *Staff) (isNew bool, err error)
	// ModifyBySelf update by self
	ModifyBySelf(uid, password string, staff *Staff) error
}

// PasswordStore Storage for Password
type PasswordStore interface {
	// Change password by self
	PasswordChange(uid, oldPassword, newPassword string) error
	// Reset password by administrator
	PasswordReset(uid, newPassword string) error
}

// Authenticator for Authenticate
type Authenticator interface {
	// Authenticate with uid and password
	Authenticate(uid, password string) (*Staff, error)
}

// GroupStore Storage for Group
type GroupStore interface {
	AllGroup() ([]Group, error)
	GetGroup(name string) (*Group, error)
	SaveGroup(*Group) error
}
