package gapi

import (
	"context"
	"database/sql"

	db "github.com/Annongkhanh/Simple_bank/db/sqlc"
	"github.com/Annongkhanh/Simple_bank/pb"
	"github.com/Annongkhanh/Simple_bank/util"
	"github.com/Annongkhanh/Simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	if violations := validateLoginUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)

	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "password not match: %s", err)

	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)

	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)

	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    server.extractMetadata(ctx).UserAgent,
		ClientIp:     server.extractMetadata(ctx).ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)

	}

	rsp := &pb.LoginUserResponse{
		User:                 convertUser(user),
		SessionId:            session.ID.String(),
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpireAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpireAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
