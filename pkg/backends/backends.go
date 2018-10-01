package backends

import (
	"fmt"
	"log"

	. "github.com/wealthworks/go-debug"

	"github.com/liut/staffio/pkg/models"
)

var (
	debug = Debug("staffio:backends")
)

// save staff
func (s *serviceImpl) SaveStaff(staff *models.Staff) error {
	if staff.EmployeeNumber < 1 {
		newId, err := NextStaffID()
		if err != nil {
			return err
		}
		staff.EmployeeNumber = newId
	}
	isNew, err := s.Save(staff)
	if err == nil {
		if isNew {
			log.Printf("new staff %v", staff)
			err = s.passwordForgotPrepare(staff)
			if err != nil {
				log.Printf("email of new user password send ERR %s", err)
			} else {
				log.Print("send email OK")
			}
		}
	} else {
		log.Printf("SaveStaff %s ERR %s", staff.Uid, err)
	}
	return err
}

func (s *serviceImpl) InGroup(gname, uid string) bool {
	g, err := s.GetGroup(gname)
	if err != nil {
		log.Printf("GetGroup %s ERR %s", gname, err)
		return false
	}
	// log.Printf("check uid %s in %v", uid, g)
	return g.Has(uid)
}

func (s *serviceImpl) ProfileModify(uid, password string, staff *models.Staff) error {
	if uid != staff.Uid {
		return fmt.Errorf("mismatch uid %s and %s", uid, staff.Uid)
	}
	return s.ModifyBySelf(uid, password, staff)
}

// 返回下一个员工ID
func NextStaffID() (eid int, err error) {
	err = withDbQuery(func(db dber) error {
		return db.Get(&eid, "SELECT nextval('staff_id_seq')")
	})
	return
}

func WriteUserLog(uid, subject, message string) error {
	qs := func(db dber) error {
		_, err := db.Exec("INSERT INTO user_log(uid, subject, body) VALUES($1, $2, $3)", uid, subject, message)
		return err
	}
	return withDbQuery(qs)
}
