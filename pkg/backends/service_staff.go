package backends

import (
	"context"
	"fmt"

	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/team"
)

// save staff
func (s *serviceImpl) SaveStaff(ctx context.Context, staff *models.Staff) error {
	if len(staff.EmployeeNumber) == 0 {
		newID, err := NextStaffID()
		if err != nil {
			return err
		}
		staff.EmployeeNumber = fmt.Sprintf("%04d", newID)
	}
	isNew, err := s.Save(staff)
	if err == nil {
		if isNew {
			logger().Infow("net staff", "staff", staff)
			err = s.passwordForgotPrepare(SiteFromContext(ctx), staff)
			if err != nil {
				logger().Infow("email of new user password send fail", "err", err)
			} else {
				logger().Infow("send email OK")
			}
		}
	} else {
		logger().Warnw("save staff fail", "staff", staff, "err", err)
	}
	return err
}

func (s *serviceImpl) InGroup(gname, uid string) bool {
	g, err := s.GetGroup(gname)
	if err != nil {
		logger().Infow("get group fail", "name", gname, "err", err)
		return false
	}
	// log.Printf("check uid %s in %v", uid, g)
	return g.Has(uid)
}

func (s *serviceImpl) InGroupAny(uid string, names ...string) bool {
	for _, gn := range names {
		g, err := s.GetGroup(gn)
		if err == nil && g.Has(uid) {
			return true
		}
	}
	logger().Infow("inGroupAny fail", "uid", uid, "groups", names)
	return false
}

func (s *serviceImpl) ProfileModify(uid, password string, staff *models.Staff) error {
	if uid != staff.UID {
		return fmt.Errorf("mismatch uid %s and %s", uid, staff.UID)
	}
	return s.ModifyBySelf(uid, password, staff)
}

// NextStaffID 返回下一个员工ID
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

// StoreTeamAndStaffs = svc.SaveTeamAndStaffs
func StoreTeamAndStaffs(ctx context.Context, svc Servicer, team *team.Team, staffs models.Staffs) (err error) {
	return svc.SaveTeamAndStaffs(ctx, team, staffs)
}

// SaveTeamAndStaffs ...
func (s *serviceImpl) SaveTeamAndStaffs(ctx context.Context, team *team.Team, staffs models.Staffs) (err error) {

	for _, staff := range staffs {
		if err = s.SaveStaff(ctx, &staff); err != nil {
			logger().Infow("save staff fail", "staff", staff, "err", err)
			return
		}
		logger().Debugw("bulk save staff ok", "cn", staff.CommonName, "uid", staff.UID)
	}

	err = s.Team().Store(team)
	if err != nil {
		logger().Infow("bulk save team fail", "name", team.Name, "err", err)
		return
	}

	logger().Infow("bulk team saved OK", "name", team.Name, "leaders", team.Leaders)
	for _, leader := range team.Leaders {
		err = s.Team().AddManager(team.ID, leader)
		if err != nil {
			logger().Infow("bulk add manager fail", "leader", leader, "teamID", team.ID, "err", err)
		}
	}

	return
}
