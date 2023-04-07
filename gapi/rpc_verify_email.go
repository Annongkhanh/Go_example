package gapi

import (
	"context"

	db "github.com/Annongkhanh/Simple_bank/db/sqlc"
	"github.com/Annongkhanh/Simple_bank/pb"
	"github.com/Annongkhanh/Simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {


	if violations := validateVerifyEmailRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateVerifyEmailParams{
		ID: req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	} 


	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email: %s", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {

		if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
			violations = append(violations, fieldViolation("email_id", err))
		}

		if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
			violations = append(violations, fieldViolation("secret_code", err))
		}
	return violations
}
