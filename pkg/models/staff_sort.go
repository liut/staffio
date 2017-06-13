package models

import (
	"sort"
)

var (
	ByUid = By(func(p1, p2 *Staff) bool {
		return p1.Uid < p2.Uid
	})
)

// By is the type of a "less" function that defines the ordering of its Staff arguments.
type By func(p1, p2 *Staff) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(staffs []*Staff) {
	ps := &staffSorter{
		staffs: staffs,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type staffSorter struct {
	staffs []*Staff
	by     func(p1, p2 *Staff) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *staffSorter) Len() int {
	return len(s.staffs)
}

// Swap is part of sort.Interface.
func (s *staffSorter) Swap(i, j int) {
	s.staffs[i], s.staffs[j] = s.staffs[j], s.staffs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *staffSorter) Less(i, j int) bool {
	return s.by(s.staffs[i], s.staffs[j])
}
