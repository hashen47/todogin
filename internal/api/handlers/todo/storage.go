package todo

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

func (s *Storage) InsertTodo(title, content string, userId int) error {
	stmt, err := s.Database.Conn.Prepare("insert into todos(title, content, user_id) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTodos(userId, limit, offset int) (*[]Todo, error) {
	stmt, err := s.Database.Conn.Prepare("select * from todos where user_id=? limit ? offset ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, limit, offset) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Content, &todo.UserId, &todo.Done); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &todos, nil
}

func (s *Storage) GetTodoById(id, userId int) (*Todo, error) {
	var stmt *sql.Stmt
	var err error

	if userId == 0 {
		stmt, err = s.Database.Conn.Prepare("select * from todos where id=?")
	} else {
		stmt, err = s.Database.Conn.Prepare("select * from todos where id=? and user_id=?")
	}

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var todo Todo

	if userId == 0 {
		err = stmt.QueryRow(id).Scan(&todo.Id, &todo.Title, &todo.Content, &todo.UserId, &todo.Done)
	} else {
		err = stmt.QueryRow(id, userId).Scan(&todo.Id, &todo.Title, &todo.Content, &todo.UserId, &todo.Done)
	}

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *Storage) GetTotalTodoCount(userId int) (int, error) {
	stmt, err := s.Database.Conn.Prepare("select count(*) from todos where user_id=?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	todoCount := 0
	if err := stmt.QueryRow(userId).Scan(&todoCount); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}

	return todoCount, nil
}

func (s *Storage) UpdateTodo(userId int, todoId int, title string, content string, done bool) error {
	stmt, err := s.Database.Conn.Prepare("update todos set title=?, content=?, done=? where user_id=? and id=?") 
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, content, done, userId, todoId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteTodo(id, userId int) error {
	stmt, err := s.Database.Conn.Prepare("delete from todos where id=? and user_id=?") 
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, userId)
	if err != nil {
		return err
	}

	return nil
}
