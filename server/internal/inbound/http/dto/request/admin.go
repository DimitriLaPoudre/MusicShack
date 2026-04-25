package request

type LoginAdmin struct {
	Password string `json:"password" binding:"required"`
}

type ChangePasswordAdmin struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
