package domain

import "time"

type User struct {
	ID         int       `json:"id"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Gender     string    `json:"gender"`
	Birthdate  time.Time `json:"birthdate"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
}
