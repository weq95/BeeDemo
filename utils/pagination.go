package utils

import "math"

type Pagination struct {
	Page      int64   `json:"page"`       //当前页
	PageSize  int64   `json:"page_size"`  //每页条数
	Total     int64   `json:"total"`      //总条数
	PageCount int64   `json:"page_count"` //总页数
	Nums      []int64 `json:"nums"`       //分页序数
	NumsCount int64   `json:"nums_count"` //总页序数
}

//创建分页
func CreatePaging(page, pageSize, total int64) *Pagination {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 15
	}

	pageCount := math.Ceil(float64(total) / float64(pageSize))
	paging := &Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		PageCount: int64(pageCount),
		NumsCount: 7,
	}

	paging.setNums()
	return paging
}

//设置分页序数信息
func (p *Pagination) setNums() {
	p.Nums = []int64{}

	if p.PageCount == 0 {
		return
	}

	half := math.Floor(float64(p.NumsCount) / float64(2))
	begin := p.Page - int64(half)

	if begin < 1 {
		begin = 1
	}

	end := begin + p.NumsCount - 1
	if end >= p.PageCount {
		begin = p.PageCount - p.NumsCount + 1
		if begin < 1 {
			begin = 1
		}

		end = p.PageCount
	}

	for i := begin; i < end; i++ {
		p.Nums = append(p.Nums, i)
	}
}
