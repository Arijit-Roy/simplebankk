package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPasswd, err := util.HashPassword("secret")
	require.NoError(t, err)
	args := CreateUserParams{
		UserName:       util.RandomOwner(),
		HashedPassword: hashedPasswd,
		FullName:       util.RandomString(10),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.UserName, args.UserName)
	require.Equal(t, args.Email, args.Email)
	require.Equal(t, args.HashedPassword, args.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChagedAt.IsZero())

	return user

}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)

	user1, err := testQueries.GetUser(context.Background(), user.UserName)

	require.NoError(t, err)
	require.Equal(t, user1.CreatedAt, user.CreatedAt)
	require.Equal(t, user1.Email, user.Email)
}
