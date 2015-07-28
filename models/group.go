package models

type Group struct {
	Name    string
	Members []string
}

var (
	EmptyGroup = &Group{"", make([]string, 0)}
)

func (g *Group) Has(member string) bool {
	for _, m := range g.Members {
		if m == member {
			return true
		}
	}
	return false
}
