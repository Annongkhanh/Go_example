package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTranferTx(t *testing.T){
	// txName := 
	store := NewStore(testDB)

	account1 := createRandomAccount(t)

	account2 := createRandomAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance )

	n := 10
	amount := int64(1000)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++{
		go func(){
			result, err := store.transferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errors <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++{
		fmt.Println("i", i)
		err := <- errors
		require.NoError(t, err)
		// Transfer 
		result := <- results
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// FromEntry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		//ToEntry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		//FromAccount
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		//ToAccount
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		fmt.Println(">> tx:", result.FromAccount.Balance, result.ToAccount.Balance)
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := account2.Balance - result.ToAccount.Balance

		require.Equal(t, diff1, -diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //diff1 = (i+1) * amount 
		
		k := int(diff1/amount)
		require.True(t, k >= 1 && k <= n )
		require.NotContains(t, existed, k)
		existed[k] = true

	} 
		//Check updated balance
		updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedAccount1)
		require.Equal(t, account1.Balance, updatedAccount1.Balance + int64(n) * amount)

		updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)
		require.NotEmpty(t, updatedAccount2)
		require.Equal(t, account2.Balance, updatedAccount2.Balance - int64(n) * amount)

		fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance )

}