package typeslink

import "time"

type ErrorResponse struct {
	Error string `json: "error`
}

type Post struct {
	ID         int64     `json:"id"`
	UserID     int       `json:"userId"`
	Content    string    `json:"content"`
	Month      int       `json:"month"`
	Breed      string    `json:"breed"`
	Gender     string    `json:"gender"`
	Vaccinated string    `json:"vaccinated"`
	CreatedAt  time.Time `json:"createdAt"`
}

/*
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
*/

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
