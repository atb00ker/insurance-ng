package models

// Users model is table that stores the basic information about the users
type Users struct {
	ID      string `json:"id" gorm:"PRIMARY_KEY"`
	IsAdmin bool   `json:"is_admin"`
}
