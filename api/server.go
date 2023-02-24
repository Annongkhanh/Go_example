package api

import (
	"fmt"
	"log"

	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/token"
	"github.com/Annongkhanh/Go_example/util"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

func NewServer(store db.Store) (*Server, error) {

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Can not load key config: ", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil{
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker}
	router := gin.Default()
	server.router = router


	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.POST("/accounts/update/", server.updateAccount)

	router.POST("/transfer", server.createTransfer)

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
