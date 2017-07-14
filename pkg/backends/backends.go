package backends

import (
	"fmt"
	"log"

	. "github.com/tj/go-debug"

	"lcgc/platform/staffio/pkg/backends/ldap"
	"lcgc/platform/staffio/pkg/models"
)

var (
	backendReady bool
	debug        = Debug("staffio:backends")

	service *Service // deprecated
)

// Prepare Deprecated
func Prepare() {
	if backendReady {
		return
	}
	service = NewService()

	backendReady = true
}

func CloseAll() {
	// closeDb()
	service.Close()
}

func LoadStaffs() []*models.Staff {
	return service.StaffStore.All()
}

func AllGroup() []models.Group {
	return service.GroupStore.AllGroup()
}

func GetGroup(name string) (*models.Group, error) {
	return service.GroupStore.GetGroup(name)
}

func GetStaff(uid string) (*models.Staff, error) {
	staff, err := service.StaffStore.Get(uid)
	if err != nil {
		log.Printf("ldap get staff with %q ERR: %s", uid, err)
		return nil, err
	}
	return staff, nil
}

// save staff
func StoreStaff(staff *models.Staff) error {
	isNew, err := service.StaffStore.Save(staff)
	if err == nil {
		if isNew {
			log.Printf("new staff %v", staff)
			err = passwordForgotPrepare(staff)
			if err != nil {
				log.Printf("email of new user password send ERR %s", err)
			} else {
				log.Print("send email OK")
			}
		}
	} else {
		log.Printf("StoreStaff %s ERR %s", staff.Uid, err)
	}
	return err
}

func DeleteStaff(uid string) error {
	err := service.StaffStore.Delete(uid)
	if err == nil {
		log.Printf("deleted uid %s", uid)
	} else {
		log.Printf("deleted uid %s ERR %s", uid, err)
	}
	return err
}

func Authenticate(uid, password string) bool {
	err := service.Authenticator.Authenticate(uid, password)
	if err != nil {
		log.Printf("Authen failed for %s, reason: %s", uid, err)
		return false
	}
	debug("%s authenticate OK", uid)
	return true
}

func PasswordChange(uid, passwordOld, passwordNew string) error {
	err := service.PasswordStore.PasswordChange(uid, passwordOld, passwordNew)
	if err != nil {
		if err == ldap.ErrLogin {
			return err
			// TODO:
		}
	}

	return err
}

func InGroup(group, uid string) bool {
	g, err := GetGroup(group)
	if err != nil {
		log.Printf("GetGroup %s ERR %s", group, err)
		return false
	}
	return g.Has(uid)
}

func ProfileModify(uid, password string, staff *models.Staff) error {
	if uid != staff.Uid {
		return fmt.Errorf("mismatch uid %s and %s", uid, staff.Uid)
	}
	return service.StaffStore.ModifyBySelf(uid, password, staff)
}

func WriteUserLog(uid, subject, message string) error {
	qs := func(db dber) error {
		_, err := db.Exec("INSERT INTO user_log(uid, subject, body) VALUES($1, $2, $3)", uid, subject, message)
		return err
	}
	return withDbQuery(qs)
}
