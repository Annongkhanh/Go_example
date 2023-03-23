package api

import (
	"fmt"

	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/token"
	"github.com/Annongkhanh/Go_example/util"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	server.setUpRouter()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	// authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	// authRoutes.POST("/accounts/update/", server.updateAccount)

	authRoutes.POST("/transfer", server.createTransfer)

	authRoutes.GET("/users/:username", server.getUser)

	server.router = router

}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
