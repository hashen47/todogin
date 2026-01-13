package api

import (
	"net/http"
	"todogin/internal/config"
	"github.com/gin-gonic/gin"
	"todogin/internal/database"
	"todogin/internal/api/handlers"
	"todogin/internal/api/handlers/todo"
	"todogin/internal/api/handlers/auth"
)

type Api struct {
	db     *database.Database 
	config *config.Config
	router *gin.Engine 
}

func ApiInit(db *database.Database, conf *config.Config) *Api {
	api := new(Api)

	api.db     = db
	api.router = gin.Default()
	api.config = conf

	// logger
	api.RegisterV1Routes()

	return api
}

func (api *Api) RegisterV1Routes() {
	v1Router := api.router.Group("/v1") 
	v1Router.Use(api.InitMiddleware())

	// health check
	v1Router.GET("/health", func(c *gin.Context) {
		errs := make(handlers.ErrsMap, 0)
		resp := handlers.NewResp(
			handlers.OK,
			map[string]any{
				"msg": "running",
			},
			nil,
			errs,
		)
		c.JSON(http.StatusOK, resp)
	})

	// auth router
	authRouter := v1Router.Group("auth") 
	auth.RegisterHandlers(authRouter)

	// todo routes
	todoRouter := v1Router.Group("todo") 
	todoRouter.Use(auth.AuthMiddleware())
	todo.RegisterHandlers(todoRouter)
}

func (api *Api) InitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("database", api.db)
        c.Set("config", api.config)
        c.Next()
    }
}

func (api *Api) Run() {
	api.router.Run(api.config.ServerAddr)
}
