package querybuilder

import "gorm.io/gorm"

// QueryParams รวม params ทั้งหมดที่รับจาก API
type QueryParams struct {
	Filters    []Filter    `json:"filters"`
	GroupBy    GroupBy     `json:"group_by"`
	Aggregates []Aggregate `json:"aggregates"`
	OrderBy    []OrderBy   `json:"order_by"`
	Pagination *Pagination `json:"pagination"`
}

// Apply นำ params ทั้งหมดไปใช้กับ GORM query
func Apply(db *gorm.DB, params QueryParams) *gorm.DB {
	db = applyFilters(db, params.Filters)
	db = applyGroupBy(db, params.GroupBy)
	db = applyAggregates(db, params.Aggregates)
	db = applyOrderBy(db, params.OrderBy)
	db = applyPagination(db, params.Pagination)
	return db
}

// ApplyWithCount Apply พร้อม count total สำหรับ pagination
func ApplyWithCount(db *gorm.DB, params QueryParams, dest any) (PageResult, error) {
	var total int64

	// count ก่อน paginate
	countDB := applyFilters(db, params.Filters)
	countDB = applyGroupBy(countDB, params.GroupBy)
	if err := countDB.Count(&total).Error; err != nil {
		return PageResult{}, err
	}

	// query จริง
	queryDB := Apply(db, params)
	if err := queryDB.Find(dest).Error; err != nil {
		return PageResult{}, err
	}

	p := params.Pagination
	if p == nil {
		p = &Pagination{Page: 1, PageSize: int(total)}
	}
	p.normalize()

	return PageResult{
		Data:       dest,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: p.totalPages(total),
	}, nil
}
