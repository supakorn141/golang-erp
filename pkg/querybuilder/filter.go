package querybuilder

import "strings"

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

// Filter กำหนดเงื่อนไขกรอง 1 รายการ
type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    any            `json:"value"`
}

func sanitizeField(field string) string {
	field = strings.TrimSpace(field)
	field = strings.ReplaceAll(field, ";", "")
	field = strings.ReplaceAll(field, "'", "")
	field = strings.ReplaceAll(field, "\"", "")
	return field
}
