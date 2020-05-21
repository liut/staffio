package weekly

import (
	"fmt"
)

type SortField struct {
	Field   string `json:"field" `
	Reverse bool   `json:"reverse" `
}

type ListSort []*SortField

type ListPager struct {
	Size   int `json:"size" `
	Offset int `json:"offset" `
}

func (this *ListPager) Sql() string {
	if this.Size < 1 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", this.Size, this.Offset)
}

func (this *ListSort) Check(fields []string) bool {
	if this == nil {
		return true
	}
	for _, sort := range *this {
		pass := false
		for _, field := range fields {
			if field == sort.Field {
				pass = true
			}
		}
		if !pass {
			return false
		}
	}
	return true
}

func (this *ListSort) Sql() string {
	if this == nil {
		return ""
	}
	sql := ""
	for _, sort := range *this {
		if sql == "" {
			sql = sql + " ORDER BY "
		} else {
			sql = sql + ", "
		}
		sql = sql + sort.Field
		if sort.Reverse {
			sql = sql + " DESC"
		}
	}
	return sql
}
