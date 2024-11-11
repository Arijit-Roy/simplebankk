package api

import (
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/util"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	UserName string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		UserName:       req.UserName,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})
	resp := CreateUserResponse{
		UserName: user.UserName,
		FullName: user.FullName,
		Email:    user.Email,
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
