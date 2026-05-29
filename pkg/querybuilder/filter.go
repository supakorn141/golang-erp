package querybuilder

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FilterOperator string

const (
	Eq        FilterOperator = "eq"
	Ne        FilterOperator = "ne"
	Gt        FilterOperator = "gt"
	Gte       FilterOperator = "gte"
	Lt        FilterOperator = "lt"
	Lte       FilterOperator = "lte"
	Like      FilterOperator = "like"
	In        FilterOperator = "in"
	NotIn     FilterOperator = "not_in"
	IsNull    FilterOperator = "is_null"
	IsNotNull FilterOperator = "is_not_null"
)

// Filter กำหนดเงื่อนไขกรองข้อมูล 1 รายการ
type Filter struct {
	Field    string         `json:"field"`    // ชื่อ column
	Operator FilterOperator `json:"operator"` // eq, ne, gt, gte, lt, lte, like, in, not_in, is_null, is_not_null
	Value    any            `json:"value"`    // ค่าที่ใช้เปรียบเทียบ (is_null/is_not_null ไม่ต้องส่ง)
}

func applyFilters(db *gorm.DB, filters []Filter) *gorm.DB {
	for _, f := range filters {
		col := sanitizeField(f.Field)
		switch f.Operator {
		case Eq:
			db = db.Where(fmt.Sprintf("%s = ?", col), f.Value)
		case Ne:
			db = db.Where(fmt.Sprintf("%s != ?", col), f.Value)
		case Gt:
			db = db.Where(fmt.Sprintf("%s > ?", col), f.Value)
		case Gte:
			db = db.Where(fmt.Sprintf("%s >= ?", col), f.Value)
		case Lt:
			db = db.Where(fmt.Sprintf("%s < ?", col), f.Value)
		case Lte:
			db = db.Where(fmt.Sprintf("%s <= ?", col), f.Value)
		case Like:
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), fmt.Sprintf("%%%v%%", f.Value))
		case In:
			db = db.Where(fmt.Sprintf("%s IN ?", col), f.Value)
		case NotIn:
			db = db.Where(fmt.Sprintf("%s NOT IN ?", col), f.Value)
		case IsNull:
			db = db.Where(fmt.Sprintf("%s IS NULL", col))
		case IsNotNull:
			db = db.Where(fmt.Sprintf("%s IS NOT NULL", col))
		}
	}
	return db
}

// sanitizeField ป้องกัน SQL injection จากชื่อ field
func sanitizeField(field string) string {
	field = strings.TrimSpace(field)
	field = strings.ReplaceAll(field, ";", "")
	field = strings.ReplaceAll(field, "'", "")
	field = strings.ReplaceAll(field, "\"", "")
	return field
}
