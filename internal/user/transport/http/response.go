package user

type UserResponse struct {
	Id       string  `json:"id"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Username string  `json:"username"`
}
