package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}

}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
	}


	return tx.Commit()
}


type TransferTxParams struct{
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

// Add new transfer record, add account entries and update accounts's balance
func (store *Store) transferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		// Update balance 
		account1, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		argAccount1 := AddAccountBalanceParams{
			ID: account1.ID,
			Amount: -arg.Amount,
		}
		result.FromAccount, err = q.AddAccountBalance(ctx, argAccount1)
		if err != nil {
			return err
		}
		// fmt.Println("account1 balance:", account1.Balance)

		account2, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		argAccount2 := AddAccountBalanceParams{
			ID: account2.ID,
			Amount:  arg.Amount,
		}
		result.ToAccount, err = q.AddAccountBalance(ctx, argAccount2)
		if err != nil {
			return err
		}



		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{

		// })
		return nil
	})

	return result, err

}