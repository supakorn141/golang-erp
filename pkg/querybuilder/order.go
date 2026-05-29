package querybuilder

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

// OrderBy กำหนดการเรียงลำดับ
type OrderBy struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}
