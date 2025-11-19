package request

type RegisterRequest struct {
	Account  string `json:"account" binding:"required,min=3,max=20"`
	Nickname string `json:"nickname" binding:"required,min=1,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}
