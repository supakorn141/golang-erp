package querybuilder

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AggregateFunction string

const (
	FuncSum   AggregateFunction = "sum"
	FuncCount AggregateFunction = "count"
	FuncAvg   AggregateFunction = "avg"
	FuncMin   AggregateFunction = "min"
	FuncMax   AggregateFunction = "max"
)

// Aggregate กำหนดการคำนวณค่าสรุป
type Aggregate struct {
	Function AggregateFunction `json:"function"` // sum, count, avg, min, max
	Field    string            `json:"field"`    // column ที่ต้องการคำนวณ (* สำหรับ count)
	Alias    string            `json:"alias"`    // ชื่อผลลัพธ์ใน response
}

func applyAggregates(db *gorm.DB, aggregates []Aggregate) *gorm.DB {
	if len(aggregates) == 0 {
		return db
	}
	parts := make([]string, 0, len(aggregates)+1)
	parts = append(parts, "*")
	for _, a := range aggregates {
		fn := strings.ToUpper(string(a.Function))
		col := sanitizeField(a.Field)
		alias := sanitizeField(a.Alias)
		if alias == "" {
			alias = fmt.Sprintf("%s_%s", strings.ToLower(fn), col)
		}
		parts = append(parts, fmt.Sprintf("%s(%s) AS %s", fn, col, alias))
	}
	return db.Select(strings.Join(parts, ", "))
}
