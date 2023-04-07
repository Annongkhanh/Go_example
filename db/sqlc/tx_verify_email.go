package db

import (
	"context"
	"database/sql"
)

type VerifyEmailTxResult struct { 
	User User
	VerifyEmail VerifyEmail
}


func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID: arg.ID,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		return err

	})

	return result, err

}
