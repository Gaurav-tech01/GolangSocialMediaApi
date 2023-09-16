package models

type User struct {
	First_name      string   `json:"first_name" validate:"required,min=2,max=50"`
	Last_name       string   `json:"last_name" validate:"required,min=2,max=50"`
	Email           string   `json:"email" validate:"required,max=50"`
	Password        string   `json:"password" validate:"required,min=5"`
	Profile_picture string   `json:"profile_picture"`
	Friends         []string `json:"friends"`
	Location        string   `json:"location"`
	Occupation      string   `json:"occupation"`
}
