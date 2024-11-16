package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cant create token %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()
	router.POST("/accounts/transfer", server.TransferAmount)
	router.POST("/accounts", server.createAccount)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/account/:id", server.DeleteAccount)

	server.router = router

	return server, nil

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
