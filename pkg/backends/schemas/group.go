package schemas

// Group ...
type Group struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

// vars
var (
	EmptyGroup = &Group{"", "", make([]string, 0)}
)

// Has
func (g *Group) Has(member string) bool {
	for _, m := range g.Members {
		if m == member {
			return true
		}
	}
	return false
}

// func (g *Group) GetName() string {
// 	return g.Name
// }

// func (g *Group) GetDescription() string {
// 	return g.Description
// }

// func (g *Group) GetMembers() []string {
// 	return g.Members
// }
