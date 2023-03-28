package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Annongkhanh/Simple_bank/db/sqlc"
	"github.com/Annongkhanh/Simple_bank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	if arg.FromAccountID == arg.ToAccountID {
		error := errors.New("can not transfer money to the same account")
		ctx.JSON(http.StatusBadRequest, errorResponse(error))
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(req.FromAccountID, req.Currency, ctx)
	if !valid {
		return
	}
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("from account is not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	_, valid = server.validAccount(req.ToAccountID, req.Currency, ctx)

	if !valid {
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(accountID int64, currency string, ctx *gin.Context) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency == db.Currency(currency) {
		return account, true
	} else {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false

	}

}
