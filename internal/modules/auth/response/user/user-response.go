package user

type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Admin User"`
	Email string `json:"email" example:"admin@example.com"`
}
