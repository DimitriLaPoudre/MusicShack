package request

import "github.com/DimitriLaPoudre/MusicShack/internal/model"

type LoginForm struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanumunicode|contains=_"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

func (req LoginForm) IntoLoginForm() model.LoginForm {
	return model.LoginForm{
		Username: req.Username,
		Password: req.Password,
		Remember: req.Remember,
	}
}
