package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserSignUpReq struct {
	Name     string `json:"name" binding:"required,min=5,max=255"`
	Email    string `json:"email" binding:"required,email,min=5,max=255"`
	Password string `json:"password" binding:"required,min=8,max=16"`
}

type UserSignInReq struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type UserCustomClaim struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}
