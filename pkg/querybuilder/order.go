package querybuilder

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

// OrderBy กำหนดการเรียงลำดับ
type OrderBy struct {
	Field     string        `json:"field"`     // column ที่ต้องการเรียง
	Direction SortDirection `json:"direction"` // asc หรือ desc
}

func applyOrderBy(db *gorm.DB, orders []OrderBy) *gorm.DB {
	for _, o := range orders {
		col := sanitizeField(o.Field)
		dir := "ASC"
		if strings.EqualFold(string(o.Direction), "desc") {
			dir = "DESC"
		}
		db = db.Order(fmt.Sprintf("%s %s", col, dir))
	}
	return db
}
