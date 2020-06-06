package qqexmail

import (
	"fmt"
	"strconv"
	"strings"

	"fhyx.online/tencent-api-go/exmail"

	"github.com/liut/staffio/pkg/models"
)

var (
	EmailDomain string
)

func GetStaffFromExmail(email string) (*models.Staff, error) {
	user, err := exmail.GetUser(email)
	if err != nil {
		return nil, err
	}

	logger().Debugw("got from exmail", "user", user)
	sn, gn := models.SplitName(user.Name)

	// log.Printf("got %q %q %q", user.Name, sn, gn)

	staff := &models.Staff{
		UID:          strings.Split(user.Alias, "@")[0],
		Email:        user.Alias,
		CommonName:   user.Name,
		Surname:      sn,
		GivenName:    gn,
		EmployeeType: user.Title,
		Mobile:       user.Mobile,
		Gender:       models.Gender(user.Gender).String(),
	}
	eid, _ := strconv.Atoi(user.ExtId)
	if eid > 0 {
		staff.EmployeeNumber = eid
	}
	return staff, nil
}

func GetEmailAddress(uid string) string {
	return fmt.Sprintf("%s@%s", uid, EmailDomain)
}

func CheckMailUnseen(uid string) int {
	email := GetEmailAddress(uid)
	count, err := exmail.CountNewMail(email)
	if err != nil {
		logger().Infow("check mail fail", "uid", uid, "err", err)
	}
	return count
}

func GetMailEntryUrl(uid string) string {
	email := GetEmailAddress(uid)
	str, err := exmail.GetLoginURL(email)
	if err != nil {
		logger().Infow("get login url fail", "uid", uid, "err", err)
	}
	return str
}
