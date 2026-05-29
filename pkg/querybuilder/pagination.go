package querybuilder

// Pagination กำหนดการแบ่งหน้า
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PageResult ผลลัพธ์พร้อม metadata การแบ่งหน้า
type PageResult struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// QueryParams รวม params ทั้งหมดที่รับจาก API
type QueryParams struct {
	Filters    []Filter    `json:"filters"`
	GroupBy    GroupBy     `json:"group_by"`
	Aggregates []Aggregate `json:"aggregates"`
	OrderBy    []OrderBy   `json:"order_by"`
	Pagination *Pagination `json:"pagination"`
}

func (p *Pagination) normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *Pagination) offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) TotalPages(total int64) int {
	pages := int(total) / p.PageSize
	if int(total)%p.PageSize > 0 {
		pages++
	}
	return pages
}
