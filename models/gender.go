package models

type Gender uint8

const (
	Unknown Gender = 0 + iota
	Male
	Female
)

var genderKeys = []string{"Unknown", "Male", "Famale"}

func (this Gender) String() string {
	return genderKeys[this]
}
