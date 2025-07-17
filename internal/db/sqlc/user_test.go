package db

import (
	"clove/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomPasswordHash())
	require.NoError(t, err)
	arg := CreateUserParams{
		FullName:     util.RandomFullName(),
		Email:        util.RandomGmail(),
		PhoneNumber:  util.RandomPhoneNumber(),
		PasswordHash: hashedPassword,
		Role:         util.RandomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Role, user.Role)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Role, user2.Role)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}
