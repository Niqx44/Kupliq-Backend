package models

type User struct {
	ID       int    `json:"id"`
	RoleID   int    `json:"id_role"`
	Username string `json:"username"` // Bisa username atau email
	Password string `json:"password"`
}
