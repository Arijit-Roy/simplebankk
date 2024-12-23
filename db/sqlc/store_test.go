package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(dbConn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> before", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	// result, err := store.TransferTx(context.Background(), TransferTxParams{
	// 	FromAccountID: account1.ID,
	// 	ToAccountID:   account2.ID,
	// 	Amount:        amount,
	// })

	// require.NoError(t, err)
	// require.Equal(t, result.Transfer.FromAccountID, account1.ID)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx-%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
		result := <-results
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">>tz", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		fmt.Println(">>>>>>>>k", k)

		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	// require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-(int64(n)*amount))
	// require.Equal(t, updatedAccount2.Balance, account2.Balance+(int64(n)*amount))

	fmt.Println(">> after", account1.Balance, account2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(dbConn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> before", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	// result, err := store.TransferTx(context.Background(), TransferTxParams{
	// 	FromAccountID: account1.ID,
	// 	ToAccountID:   account2.ID,
	// 	Amount:        amount,
	// })

	// require.NoError(t, err)
	// require.Equal(t, result.Transfer.FromAccountID, account1.ID)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx-%d", i+1)
		fromAccountID := account1.ID
		ToAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			ToAccountID = account1.ID

		}

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}
	// existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

	}

	// check final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

	fmt.Println(">> after", account1.Balance, account2.Balance)

}
