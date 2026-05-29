package querybuilder

// GroupBy กำหนด field ที่ต้องการ group
type GroupBy struct {
	Fields []string `json:"fields"`
}
