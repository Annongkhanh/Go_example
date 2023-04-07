package gapi

import (
	"context"
	"database/sql"

	db "github.com/Annongkhanh/Simple_bank/db/sqlc"
	pb "github.com/Annongkhanh/Simple_bank/pb"
	"github.com/Annongkhanh/Simple_bank/util"
	"github.com/Annongkhanh/Simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's information")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		Fullname: sql.NullString{
			String: req.GetFullname(),
			Valid:  req.Fullname != nil},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)

		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user %s not found", req.GetUsername())
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.Fullname != nil {
		if err := val.ValidateFullname(req.GetFullname()); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
