package data

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := &Account{
		FirstName:    "Miro",
		LastName:     "Tama",
		Email:        "techno@mail.com",
		Password:     "passport1234",
		UserType:     "ADMIN",
		Avatar:       "default.png",
		Uuid:         "0f255a0d-a6b0-4fe0-a147-1ab1035a5f5a",
		Token:        `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJGaXJzdE5hbWUiOiJtaXJvc2xhdiIsIkxhc3ROYW1lIjoicGFudG9zIiwiRW1haWwiOiJ0ZWNobm9AZ21haWwuY29tIiwiVXNlclR5cGUiOiJVU0VSIiwiVXVpZCI6IjBmMjU1YTBkLWE2YjAtNGZlMC1hMTQ3LTFhYjEwMzVhNWY1YSIsImV4cCI6MTY3MDU5MjAzNX0.zcMxxA88Gwz_mucwlLBAd7fUXfjhj8HP0ulHyIYHFbU`,
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJGaXJzdE5hbWUiOiIiLCJMYXN0TmFtZSI6IiIsIkVtYWlsIjoiIiwiVXNlclR5cGUiOiIiLCJVdWlkIjoiIiwiZXhwIjoxNjcxMTEwNDM1fQ.ogUttNj2IlZakLiieIEr8QX4JDkOj2GicdCQq6L8UP4`,
	}

	err := testQueries.CreateAccout(arg)
	require.NoError(t, err)
}
