package backends

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wealthworks/go-tencent-api/exmail"
	"github.com/wealthworks/go-tencent-api/exwechat"

	"github.com/liut/staffio/pkg/models"
)

var (
	EmailDomain string
)

func GetStaffFromWechatUser(user *exwechat.User) *models.Staff {
	logger().Debugw("got from exmail", "user", user)
	sn, gn := models.SplitName(user.Name)
	staff := &models.Staff{
		Uid:          user.UID,
		Email:        user.Email,
		CommonName:   user.Name,
		Surname:      sn,
		GivenName:    gn,
		EmployeeType: user.Title,
		Mobile:       user.Mobile,
		Gender:       models.Gender(user.Gender),
	}
	return staff
}

func GetStaffFromExmail(email string) (*models.Staff, error) {
	user, err := exmail.GetUser(email)
	if err != nil {
		return nil, err
	}

	logger().Debugw("got from exmail", "user", user)
	sn, gn := models.SplitName(user.Name)

	// log.Printf("got %q %q %q", user.Name, sn, gn)
	eid, _ := strconv.Atoi(user.ExtId)

	staff := &models.Staff{
		Uid:          strings.Split(user.Alias, "@")[0],
		Email:        user.Alias,
		CommonName:   user.Name,
		Surname:      sn,
		GivenName:    gn,
		EmployeeType: user.Title,
		Mobile:       user.Mobile,
		Gender:       models.Gender(user.Gender),
	}
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
