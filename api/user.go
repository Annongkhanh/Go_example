package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/util"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreateAt          time.Time `json:"create_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.Fullname,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return

			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResponse{
		Username:          user.Username,
		Fullname:          user.Fullname,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreateAt:          user.CreateAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResponse{
		Username:          user.Username,
		Fullname:          user.Fullname,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreateAt:          user.CreateAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
