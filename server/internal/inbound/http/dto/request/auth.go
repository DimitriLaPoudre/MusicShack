package request

type LoginForm struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanumunicode|contains=_"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}
