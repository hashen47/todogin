package todo

import (
	"log"
	"errors"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"todogin/internal/api/handlers"
)

func RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/", getTodos)
	router.POST("/create", createTodo)
	router.PUT("/update", updateTodo)
	router.DELETE("/destroy", deleteTodo)
}

func getTodos(c *gin.Context) {
	var req TodoGetReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errs, err := handlers.GetErrorMsgs(req, err) 
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			err,
			errs,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userId := c.MustGet("user_id").(int)

	storage := NewStorage(c)
	todos, err := storage.GetTodos(userId, req.Limit, req.Offset)

	errs := make(handlers.ErrsMap, 0)
	resp := handlers.NewResp(
		handlers.FAIL,
		map[string]any{},
		err,
		errs,
	)
	if err != nil {
		log.Printf("(storage.GetTodos) Err: %v\n", err)
		resp["error"] = "Internal Server Error" 
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	totalTodoCount, err := storage.GetTotalTodoCount(userId)
	if err != nil {
		log.Printf("(storage.GetTotalTodoCount) Err: %v\n", err)
		resp["error"] = "Internal Server Error" 
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = handlers.NewResp(
		handlers.OK,
		map[string]any{
			"todos": *todos,
			"total_todos_count": totalTodoCount,
		},
		nil,
		errs,
	)
	c.JSON(http.StatusOK, resp)
}

func createTodo(c *gin.Context) {
	var req TodoCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errs, err := handlers.GetErrorMsgs(req, err)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			err,
			errs,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userId := c.MustGet("user_id").(int)

	storage := NewStorage(c)
	err := storage.InsertTodo(req.Title, req.Content, userId)
	if err != nil {
		errs := make(handlers.ErrsMap, 0)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			nil,
			errs,
		)
		log.Printf("(storage.InsertTodo) Err: %v\n", err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	errs := make(handlers.ErrsMap, 0)
	resp := handlers.NewResp(
		handlers.OK,
		map[string]any{
			"msg": "todo creation success",
		},
		nil,
		errs,
	)
	c.JSON(http.StatusCreated, resp)
}

func updateTodo(c *gin.Context) {
	var req TodoUpdateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errs, err := handlers.GetErrorMsgs(req, err) 
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			err,
			errs,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userId := c.MustGet("user_id").(int)
	storage := NewStorage(c)

	errs := make(handlers.ErrsMap, 0)
	resp := handlers.NewResp(
		handlers.FAIL,
		map[string]any{},
		nil,
		errs,
	)
	_, err := storage.GetTodoById(req.Id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp["error"] = "invalid todo id, todo not found"
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		log.Printf("(storage.GetTodoById) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = storage.UpdateTodo(userId, req.Id, req.Title, req.Content, req.Done)
	if err != nil {
		log.Printf("(storage.UpdateTodo) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp = handlers.NewResp(
		handlers.OK,
		map[string]any{
			"msg": "todo is updated",
		},
		nil,
		errs,
	)
	c.JSON(http.StatusOK, resp) 
}

func deleteTodo(c *gin.Context) {
	var req TodoDeleteReq

	if err := c.ShouldBindJSON(&req); err != nil {
		errs, err := handlers.GetErrorMsgs(req, err)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			err,
			errs,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userId := c.MustGet("user_id").(int)
	storage := NewStorage(c)

	errs := make(handlers.ErrsMap, 0)
	resp := handlers.NewResp(
		handlers.FAIL,
		map[string]any{},
		nil,
		errs,
	)

	_, err := storage.GetTodoById(req.Id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp["error"] = "invalid todo id, todo not found" 
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		log.Printf("(storage.GetTodoById) Err: %v\n", err)
		resp["error"] = "Internal Server Error" 
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = storage.DeleteTodo(req.Id, userId)
	if err != nil {
		log.Printf("(storage.DeleteTodo) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp = handlers.NewResp(
		handlers.OK,
		map[string]any{
			"msg": "todo has deleted",
		},
		nil,
		errs,
	)
	c.JSON(http.StatusOK, resp)
}
