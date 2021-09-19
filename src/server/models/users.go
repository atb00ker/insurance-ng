package models

type Users struct {
	Id      string `json:"id" gorm:"PRIMARY_KEY"`
	IsAdmin bool   `json:"is_admin"`
}
