package querybuilder

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
	Field    string            `json:"field"`    // column ที่คำนวณ
	Alias    string            `json:"alias"`    // ชื่อ column ใน result
}
