package utils

import (
	"net/url"
	"strconv"
)

type Pagination struct {
	PageSize int
	PageNo   int
}

func ParsePagination(values url.Values) *Pagination {
	p := Pagination{}
	{
		v, ok := values["page_size"]
		if !ok || len(v) == 0 {
			return nil
		}
		size, err := strconv.Atoi(v[0])
		if err != nil {
			return nil
		}
		p.PageSize = size
	}

	{
		v, ok := values["page_no"]
		if !ok || len(v) == 0 {
			return nil
		}
		no, err := strconv.Atoi(v[0])
		if err != nil {
			return nil
		}
		p.PageNo = no
	}

	if !p.Verify() {
		return nil
	}

	return &p
}

func (p *Pagination) Verify() bool {
	return p.PageSize > 0 && p.PageNo > 0
}

func (p *Pagination) LimitOffset() (uint64, uint64) {
	return uint64(p.PageSize), uint64(p.PageSize * (p.PageNo - 1))
}
