package restful

import (
	"ptcg_trader/internal/errors"
)

//
// define common usage struct
//

// DefaultPaginationPerPage ...
const DefaultPaginationPerPage int = 50

// PaginationQuery ...
type PaginationQuery struct {
	Page    int `query:"page" json:"page" example:"1"`
	PerPage int `query:"per_page" json:"per_page" example:"50"`
}

// ValidateAndSet validate input and set the default page or per_page values
func (p *PaginationQuery) ValidateAndSet() error {
	if p.Page < 0 || p.PerPage < 0 {
		return errors.Wrap(errors.ErrBadRequest, "param Page or PerPage invalid")
	}

	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = DefaultPaginationPerPage
	}

	return nil
}

// Pagination define resource metadata
type Pagination struct {
	Page       int `json:"page" example:"1"`
	PerPage    int `json:"per_page" example:"50"`
	TotalCount int `json:"total_count" example:"80"`
	TotalPage  int `json:"total_page" example:"2"`
}

// NewPagination new pagination
func NewPagination(page, perPage, total int) Pagination {
	totalPage := total / perPage
	if total%perPage > 0 {
		totalPage++
	}

	return Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		TotalPage:  totalPage,
	}
}

//
// define known errors pattens for swagger
//

// StautsBadRequestResp ...
type StautsBadRequestResp struct {
	Code    string `json:"code" example:"400000"`
	Message string `json:"message" example:"Invailed Input"`
}

// StatusUnauthorizedResp ...
type StatusUnauthorizedResp struct {
	Code    string `json:"code" example:"401000"`
	Message string `json:"message" example:"Unauthorized"`
}

// StatusForbiddenResp ...
type StatusForbiddenResp struct {
	Code    string `json:"code" example:"403001"`
	Message string `json:"message" example:"Forbidden"`
}

// StautsNotFoundResp ...
type StautsNotFoundResp struct {
	Code    string `json:"code" example:"404001"`
	Message string `json:"message" example:"Resource Not Found"`
}

// StatusInternalServerErrorResp ...
type StatusInternalServerErrorResp struct {
	Code    string `json:"code" example:"500000"`
	Message string `json:"message" example:"Internal Server Error"`
}
