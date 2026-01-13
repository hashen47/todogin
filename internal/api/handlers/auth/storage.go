package auth

import (
	"errors"
	"database/sql"
	"github.com/gin-gonic/gin"
	"todogin/internal/database"
)

type Storage struct {
	*database.Database
}

func NewStorage(c *gin.Context) *Storage {
	db := c.MustGet("database").(*database.Database)
	return &Storage{Database: db}
}

func (s *Storage) GetUserByEmail(email string) (*User, error) {
	stmt, err := s.Database.Conn.Prepare("select * from users where email=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Name, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetUserById(id int) (*User, error) {
	stmt, err := s.Database.Conn.Prepare("select * from users where id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) InsertUser(name, email, password string) error {
	_, err := s.GetUserByEmail(email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	stmt, err := s.Database.Conn.Prepare("insert into users(name, email, password) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email, password)
	if err != nil {
		return err
	}

	return nil
}
