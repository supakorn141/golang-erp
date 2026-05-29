package querybuilder

import "gorm.io/gorm"

// Pagination กำหนดการแบ่งหน้า
type Pagination struct {
	Page     int `json:"page"`      // หน้าที่ต้องการ (เริ่มที่ 1)
	PageSize int `json:"page_size"` // จำนวนรายการต่อหน้า (default 20, max 100)
}

// PageResult ผลลัพธ์พร้อม metadata
type PageResult struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
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

func (p *Pagination) totalPages(total int64) int {
	pages := int(total) / p.PageSize
	if int(total)%p.PageSize > 0 {
		pages++
	}
	return pages
}

func applyPagination(db *gorm.DB, p *Pagination) *gorm.DB {
	if p == nil {
		return db
	}
	p.normalize()
	return db.Offset(p.offset()).Limit(p.PageSize)
}
