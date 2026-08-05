package request

type SearchSetupQuery struct {
	Q string `form:"q" binding:"required"`
	PaginationQuery
}

type SearchQuery struct {
	Q        string `form:"q" binding:"required"`
	Provider string `form:"provider" binding:"required"`
	PaginationQuery
}
