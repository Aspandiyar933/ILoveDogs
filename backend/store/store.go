package store

import (
	"database/sql"
	"errors"

	"github.com/Aspandiyar933/Ilovedogs/typeslink"
)

type Store interface {
	CreatePost(p *typeslink.Post) (*typeslink.Post, error)
	GetPost(id string) (*typeslink.Post, error)
	GetUserByID(id string) (*typeslink.User, error)
	CreateUser(u *typeslink.User) (*typeslink.User, error)
	GetUserByEmail(email string) (*typeslink.LoginRequest, error)
}

type Storage struct {
	db *sql.DB 
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *typeslink.User) (*typeslink.User, error) {
	rows, err := s.db.Exec("INSERT INTO Human (name, username, email, password) VALUES(?, ?, ?, ?)", u.Name, u.UserName, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = id 
	return u, nil
}

func (s *Storage) CreatePost(p *typeslink.Post) (*typeslink.Post, error) {
	rows, err := s.db.Exec("INSERT INTO Post (content, month, breed, gender, vaccinated) VALUES(?, ?, ?, ?, ?)", 
	p.Content, p.Month, p.Breed, p.Gender, p.Vaccinated)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id
	return p, nil
}

func (s *Storage) GetPost(id string) (*typeslink.Post, error) {
	var p typeslink.Post
	err := s.db.QueryRow("SELECT id, userId, content, month, breed, gender, vaccinated, createAt FROM Post WHERE id = ?", id).Scan(
		&p.ID ,&p.Content, &p.Month, &p.Breed, &p.Gender, &p.Vaccinated)
	return &p, err
}



func (s *Storage) GetUserByID(id string) (*typeslink.User, error) {
	var u typeslink.User
	err := s.db.QueryRow("SELECT id, name, username, email, password FROM Human WHERE id = ?", id).Scan(
		&u.ID, &u.Name, &u.UserName, &u.Email, &u.Password)
	return &u, err
}

// GetUserByEmail retrieves a login from the database by their email address
func (s *Storage) GetUserByEmail(email string) (*typeslink.LoginRequest, error) {
    query := "SELECT id, name, username, email, password FROM users WHERE email = ?"

    row := s.db.QueryRow(query, email)

    var login *typeslink.LoginRequest

    err := row.Scan(&login.Email, &login.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("login not found")
        }
        return nil, err
    }
	login.Email = email
    return login, nil
}
