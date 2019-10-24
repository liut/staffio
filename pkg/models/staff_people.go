package models

import (
	"github.com/liut/staffio/pkg/backends/schemas"
	"github.com/liut/staffio/pkg/common"
)

type Gender = common.Gender

const (
	Unknown = common.Unknown
	Male    = common.Male
	Female  = common.Female
)

type Staff = schemas.People
type Staffs = schemas.Peoples
type Spec = schemas.Spec
