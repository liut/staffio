package models

// Storage for Staff
type StaffStore interface {
	// All browse from store, like LDAP
	All() []*Staff
	// Get with uid
	Get(uid string) (*Staff, error)
	// Delete with uid
	Delete(uid string) error
	// Save add or update
	Save(staff *Staff) (isNew bool, err error)
	// ModifyBySelf update by self
	ModifyBySelf(uid, password string, staff *Staff) error
}

// Storage for Password
type PasswordStore interface {
	// Change password by self
	PasswordChange(uid, old_password, new_password string) error
	// Reset password by administrator
	PasswordReset(uid, new_password string) error
}

// Authenticator
type Authenticator interface {
	// Authenticate with uid and password
	Authenticate(uid, password string) error
}

// Storage for Group
type GroupStore interface {
	AllGroup() []Group
	GetGroup(name string) (*Group, error)
	SaveGroup(*Group) error
}
