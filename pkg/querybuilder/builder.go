package querybuilder

import (
	"fmt"
	"strings"
)

// Builder สร้าง SQL SELECT query แบบ string ต่อกัน
type Builder struct {
	table    string
	selects  []string
	wheres   []string
	groupBys []string
	orderBys []string
	args     []any
	argIdx   int
	limitVal int
	offVal   int
}

func New(table string) *Builder {
	return &Builder{table: table, argIdx: 1}
}

// nextPlaceholder คืน $N แล้วเพิ่ม counter
func (b *Builder) nextPlaceholder(val any) string {
	ph := fmt.Sprintf("$%d", b.argIdx)
	b.args = append(b.args, val)
	b.argIdx++
	return ph
}

func (b *Builder) Select(cols ...string) *Builder {
	b.selects = append(b.selects, cols...)
	return b
}

func (b *Builder) Aggregate(fn AggregateFunction, field, alias string) *Builder {
	col := fmt.Sprintf("%s(%s) AS %s",
		strings.ToUpper(string(fn)),
		sanitizeField(field),
		sanitizeField(alias),
	)
	b.selects = append(b.selects, col)
	return b
}

func (b *Builder) Where(field string, op FilterOperator, val any) *Builder {
	col := sanitizeField(field)
	switch op {
	case Eq:
		b.wheres = append(b.wheres, fmt.Sprintf("%s = %s", col, b.nextPlaceholder(val)))
	case Ne:
		b.wheres = append(b.wheres, fmt.Sprintf("%s != %s", col, b.nextPlaceholder(val)))
	case Gt:
		b.wheres = append(b.wheres, fmt.Sprintf("%s > %s", col, b.nextPlaceholder(val)))
	case Gte:
		b.wheres = append(b.wheres, fmt.Sprintf("%s >= %s", col, b.nextPlaceholder(val)))
	case Lt:
		b.wheres = append(b.wheres, fmt.Sprintf("%s < %s", col, b.nextPlaceholder(val)))
	case Lte:
		b.wheres = append(b.wheres, fmt.Sprintf("%s <= %s", col, b.nextPlaceholder(val)))
	case Like:
		b.wheres = append(b.wheres, fmt.Sprintf("%s ILIKE %s", col, b.nextPlaceholder("%"+fmt.Sprint(val)+"%")))
	case In:
		b.wheres = append(b.wheres, fmt.Sprintf("%s = ANY(%s)", col, b.nextPlaceholder(val)))
	case NotIn:
		b.wheres = append(b.wheres, fmt.Sprintf("%s != ALL(%s)", col, b.nextPlaceholder(val)))
	case IsNull:
		b.wheres = append(b.wheres, fmt.Sprintf("%s IS NULL", col))
	case IsNotNull:
		b.wheres = append(b.wheres, fmt.Sprintf("%s IS NOT NULL", col))
	}
	return b
}

func (b *Builder) GroupBy(fields ...string) *Builder {
	for _, f := range fields {
		b.groupBys = append(b.groupBys, sanitizeField(f))
	}
	return b
}

func (b *Builder) OrderBy(field string, dir SortDirection) *Builder {
	d := "ASC"
	if strings.EqualFold(string(dir), "desc") {
		d = "DESC"
	}
	b.orderBys = append(b.orderBys, fmt.Sprintf("%s %s", sanitizeField(field), d))
	return b
}

func (b *Builder) Limit(n int) *Builder  { b.limitVal = n; return b }
func (b *Builder) Offset(n int) *Builder { b.offVal = n; return b }

// Build คืน SQL string และ args พร้อมใช้กับ db.Query
func (b *Builder) Build() (string, []any) {
	sel := "*"
	if len(b.selects) > 0 {
		sel = strings.Join(b.selects, ", ")
	}

	q := fmt.Sprintf("SELECT %s FROM %s", sel, b.table)

	if len(b.wheres) > 0 {
		q += "\nWHERE " + strings.Join(b.wheres, "\n  AND ")
	}
	if len(b.groupBys) > 0 {
		q += "\nGROUP BY " + strings.Join(b.groupBys, ", ")
	}
	if len(b.orderBys) > 0 {
		q += "\nORDER BY " + strings.Join(b.orderBys, ", ")
	}
	if b.limitVal > 0 {
		q += fmt.Sprintf("\nLIMIT %d", b.limitVal)
	}
	if b.offVal > 0 {
		q += fmt.Sprintf("\nOFFSET %d", b.offVal)
	}

	return q, b.args
}

// BuildCount คืน SQL สำหรับนับ total rows (ใช้คู่กับ pagination)
func (b *Builder) BuildCount() (string, []any) {
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", b.table)
	if len(b.wheres) > 0 {
		q += "\nWHERE " + strings.Join(b.wheres, "\n  AND ")
	}
	if len(b.groupBys) > 0 {
		q += "\nGROUP BY " + strings.Join(b.groupBys, ", ")
	}
	return q, b.args
}

// Apply รับ QueryParams แล้วใส่ทุก clause เข้า builder ให้อัตโนมัติ
func (b *Builder) Apply(params QueryParams) *Builder {
	for _, f := range params.Filters {
		b.Where(f.Field, f.Operator, f.Value)
	}
	for _, a := range params.Aggregates {
		b.Aggregate(a.Function, a.Field, a.Alias)
	}
	if len(params.GroupBy.Fields) > 0 {
		b.GroupBy(params.GroupBy.Fields...)
	}
	for _, o := range params.OrderBy {
		b.OrderBy(o.Field, o.Direction)
	}
	if params.Pagination != nil {
		p := params.Pagination
		p.normalize()
		b.Limit(p.PageSize).Offset(p.offset())
	}
	return b
}
