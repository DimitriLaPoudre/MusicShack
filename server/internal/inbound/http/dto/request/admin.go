package request

type LoginAdmin struct {
	Password string `json:"password" bindings:"required"`
}
