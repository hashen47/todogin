package auth

import (
	"log"
	"time"
	"errors"
	"strconv"
	"net/http"
	"database/sql"
	"todogin/internal/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"todogin/internal/api/handlers"
)

func RegisterHandlers(router *gin.RouterGroup) {
	router.POST("/signin", signIn)
	router.POST("/signup", signUp)
}

func signIn(c *gin.Context) {
	conf := c.MustGet("config").(*config.Config)
	lifetime, err := strconv.ParseInt(conf.JwtTokenLifetime, 10, 64)
	if err != nil {
		errs := make(handlers.ErrsMap)
		log.Printf("(strconv.ParseInt) Err: %v\n", err)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			err,
			errs,
		)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	var req UserSignInReq
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

	store := NewStorage(c)
	user, err := store.GetUserByEmail(req.Email)
	if err != nil {
		errs := make(handlers.ErrsMap, 0)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			nil,
			errs,
		)
		if errors.Is(err, sql.ErrNoRows) {
			resp["error"] = "email and password combination is wrong"
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		log.Printf("(store.GetUserByEmail) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		errs := make(handlers.ErrsMap, 0)
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			errors.New("email and password combination is wrong"),
			errs,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	claims := UserCustomClaim{
		UserId          : user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(lifetime * int64(time.Second)))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "todogin_authentication",
			Subject:   "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(conf.JwtSecretKey))
	errs := make(handlers.ErrsMap, 0)

	if err != nil {
		resp := handlers.NewResp(
			handlers.FAIL,
			map[string]any{},
			errors.New("Internal Server Error"),
			errs,
		)
		log.Printf("(token.SignedString) Err: %v\n", err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := handlers.NewResp(
		handlers.OK,
		map[string]any{
			"msg"  : "signin has completed",
			"token": ss,
		},
		nil,
		errs,
	)
	c.JSON(http.StatusOK, resp)
}

func signUp(c *gin.Context) {
	var req UserSignUpReq
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

	store := NewStorage(c)
	_, err := store.GetUserByEmail(req.Email)

	errs := make(handlers.ErrsMap, 0)
	resp := handlers.NewResp(
		handlers.FAIL,
		map[string]any{},
		nil,
		errs,
	)

	if err == nil {
		resp["error"] = "email is already taken"
		c.JSON(http.StatusBadRequest, resp)
		return
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("(store.GetUserByEmail) Err: %v\n", err)
			resp["error"] = "Internal Server Error"
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("(bcrypt.GenerateFromPassword) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = store.InsertUser(req.Name, req.Email, string(hash))
	if err != nil {
		log.Printf("(store.InsertUser) Err: %v\n", err)
		resp["error"] = "Internal Server Error"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp = handlers.NewResp(
		handlers.OK,
		map[string]any{
			"msg": "user has created",
		},
		nil,
		errs,
	)
	c.JSON(http.StatusCreated, resp)
}
