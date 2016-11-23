package backends

import (
	"log"
	"strings"

	"github.com/wealthworks/go-tencent-api/exmail"

	"lcgc/platform/staffio/models"
)

func GetStaffFromExmail(email string) (*models.Staff, error) {
	user, err := exmail.GetUser(email)
	if err != nil {
		return nil, err
	}

	sn, gn := models.SplitName(user.Name)
	log.Printf("%q %q %q", user.Name, sn, gn)

	return &models.Staff{
		Uid:            strings.Split(user.Alias, "@")[0],
		Email:          user.Alias,
		CommonName:     user.Name,
		Surname:        sn,
		GivenName:      gn,
		EmployeeNumber: user.ExtId,
		EmployeeType:   user.Title,
		Mobile:         user.Mobile,
		Gender:         models.Gender(user.Gender),
	}, nil
}
