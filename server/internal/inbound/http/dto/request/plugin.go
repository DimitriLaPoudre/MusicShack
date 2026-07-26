package request

type SearchQuery struct {
	Q string `form:"q" binding:"required"`
	PaginationQuery
}
