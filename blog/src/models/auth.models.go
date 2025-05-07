package models

//AuthForm : structure used to grab user data from /contact POST requests
type AuthForm struct {
	Email     string `json:"email" form:"email" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	UserToken string `json:"userToken"`
}
