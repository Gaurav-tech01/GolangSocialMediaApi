package models

type Post struct {
	UserId       string          `json:"userId" validate:"required"`
	First_name   string          `json:"first_name" validate:"required,min=2,max=50"`
	Last_name    string          `json:"last_name" validate:"required,min=2,max=50"`
	Location     string          `json:"location"`
	Description  string          `json:"description"`
	Picture_path string          `json:"picture_path" validate:"required,min=5"`
	User_picture string          `json:"user_picture"`
	Likes        map[string]bool `json:"likes"`
	Comments     []string        `json:"comment"`
}
