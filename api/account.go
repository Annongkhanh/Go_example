package api

import (
	"database/sql"
	"net/http"

	db "github.com/Annongkhanh/Go_example/db/sqlc"

	"github.com/gin-gonic/gin"
)


type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
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



type getAccountRequest struct{
	id    int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount (ctx *gin.Context ){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, req.id)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, err)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}