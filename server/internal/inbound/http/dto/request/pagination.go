package request

type PaginationQuery struct {
	Offset int `form:"offset,default=0" binding:"min=0"`
	Limit  int `form:"limit,default=25" binding:"min=1,max=100"`
}
