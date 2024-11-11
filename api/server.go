package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.POST("/accounts/transfer", server.TransferAmount)
	router.POST("/accounts", server.createAccount)

	router.POST("/users", server.createUser)

	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/account/:id", server.DeleteAccount)

	server.router = router

	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
