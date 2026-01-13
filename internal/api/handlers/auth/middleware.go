package auth

import (
	"log"
	"errors"
	"strings"
	"net/http"
	"database/sql"
	"todogin/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"todogin/internal/api/handlers"
)

func AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		errs := make(handlers.ErrsMap, 0)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			nil,
			errs,
		)
		if authHeader == "" {
			resp["error"] = "Authorization header not found"
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			resp["error"] = "Invalid Authorization header"
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		conf := c.MustGet("config").(*config.Config)
		token, err := jwt.ParseWithClaims(headerParts[1], &UserCustomClaim{}, func(token *jwt.Token) (any, error) {
			return []byte(conf.JwtSecretKey), nil
		})

		if err != nil {
			resp["error"] = err.Error()
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		claims, ok := token.Claims.(*UserCustomClaim)
		if !ok {
			log.Printf("[token.Claims.(*UserCustomClaim).GetUserById] Err: %v\n", err)
			resp["error"] = "Internal Server Error"
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}

		store := NewStorage(c)
		user, err := store.GetUserById(claims.UserId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				resp["error"] = "Invalid token value" 
				c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
				return
			}
			log.Printf("(store.GetUserById) Err: %v\n", err)
			resp["error"] = "Internal Server Error"
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}

		c.Set("user_id", user.Id)
		c.Next()
	}
}
