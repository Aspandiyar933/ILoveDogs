package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB 
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL!")
	return &MySQLStorage{
		db: db,
	}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {

	if err := s.createUserTable(); err != nil {
		return nil, err
	}

	if err := s.createPostTable(); err != nil {
		return nil, err
	}

	if err := s.createCommentTable(); err != nil {
		return nil, err
	}
	return s.db, nil
}

func (s *MySQLStorage) createUserTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS Human(
			id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
			name VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			PRIMARY KEY(id));
	`)

	return err
}

func (s *MySQLStorage) createPostTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS Post(
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			userId INT UNSIGNED NOT NULL,  -- Use INT UNSIGNED for userId
			content VARCHAR(255) NOT NULL,
			month INT NOT NULL,
			breed VARCHAR(255) NOT NULL,
			gender CHAR(1) NOT NULL,
			vaccinated VARCHAR(255) NOT NULL,
			createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(id),
			FOREIGN KEY(userId) REFERENCES Human(id));  -- Reference Human(id)
	`)
	
	return err
}


func (s *MySQLStorage) createCommentTable() error {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS Comment(
            postId INT UNSIGNED NOT NULL,  -- Change postId to INT UNSIGNED
            userId INT UNSIGNED NOT NULL,  -- Use INT UNSIGNED for userId
            content VARCHAR(255),
            FOREIGN KEY(postId) REFERENCES Post(id),
            FOREIGN KEY(userId) REFERENCES Human(id));
    `)
    
    return err
}




