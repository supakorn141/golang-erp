package querybuilder

import (
	"strings"

	"gorm.io/gorm"
)

// GroupBy กำหนด field ที่ต้องการ group
type GroupBy struct {
	Fields []string `json:"fields"` // รายการ column ที่ group
}

func applyGroupBy(db *gorm.DB, group GroupBy) *gorm.DB {
	if len(group.Fields) == 0 {
		return db
	}
	sanitized := make([]string, 0, len(group.Fields))
	for _, f := range group.Fields {
		sanitized = append(sanitized, sanitizeField(f))
	}
	return db.Group(strings.Join(sanitized, ", "))
}
