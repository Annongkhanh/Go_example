package api

import (
	"net/http"
	db "github.com/Annongkhanh/Go_example/db/sqlc"

	"github.com/gin-gonic/gin"
)


type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency int64 `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount (ctx *gin.Context ){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Balance: int64(0),
		Currency: db.Currency(req.Currency),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}