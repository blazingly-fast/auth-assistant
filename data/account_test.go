package data

import (
	"testing"
	"time"

	"github.com/blazingly-fast/auth-assistant/util"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *Account {

	req := &CreateAccountRequest{
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     util.RandomEmail(),
		Password:  "passport1234",
	}

	userType := "USER"
	avatar := "default.png"

	uuid := uuid.New().String()

	token, refreshToken, _ := util.GenerateAllToken(
		req.FirstName,
		req.LastName,
		req.Email,
		userType,
		uuid,
	)
	acc := NewAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		userType,
		avatar,
		uuid,
		token,
		refreshToken,
	)

	err := testQueries.CreateAccout(acc)

	require.NoError(t, err)
	return acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountByUUID(t *testing.T) {
	randAcc := createRandomAccount(t)

	acc, err := testQueries.GetAccountByField("uuid", randAcc.Uuid)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, randAcc.Email, acc.Email)
	require.Equal(t, randAcc.Uuid, acc.Uuid)
	require.Equal(t, randAcc.Token, acc.Token)
}

func TestGetAccountByEmail(t *testing.T) {
	randAcc := createRandomAccount(t)

	acc, err := testQueries.GetAccountByField("email", randAcc.Email)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, randAcc.Email, acc.Email)
	require.Equal(t, randAcc.Uuid, acc.Uuid)
	require.Equal(t, randAcc.Token, acc.Token)
}

func TestUpdateAccount(t *testing.T) {

	randAcc := createRandomAccount(t)

	firstName := util.RandomName()
	lastName := util.RandomName()

	password := "passport1234"
	updateTime := time.Now().UTC()

	hashedPassword, _ := util.HashPassword(password)

	acc := &UpdateAccountRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     randAcc.Email,
		Password:  hashedPassword,
		UserType:  "USER",
		UpdatedOn: updateTime,
	}

	err := testQueries.UpdateAccount(acc, randAcc.Uuid)

	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {

	randAcc := createRandomAccount(t)

	err := testQueries.DeleteAccount(randAcc.Uuid)
	require.NoError(t, err)

	acc, err := testQueries.GetAccountByField("uuid", randAcc.Uuid)

	require.Error(t, err)
	require.Empty(t, acc)
}

func TestGetAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.GetAccounts(5, 1)

	require.NoError(t, err)
	require.Len(t, accounts.Accounts, 5)

	for _, account := range accounts.Accounts {
		require.NotEmpty(t, account)
	}
}
