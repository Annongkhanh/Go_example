package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(int(RandomInt(32, 8)))

	hashedPassword1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	require.NoError(t, CheckPassword(password, hashedPassword1))

	wrongPassword := RandomString(int(RandomInt(32, 8)))
	if wrongPassword == password {
		require.NoError(t, CheckPassword(wrongPassword, hashedPassword1))

	} else {
		require.EqualError(t, CheckPassword(wrongPassword, hashedPassword1), bcrypt.ErrMismatchedHashAndPassword.Error())
	}

	hashedPassword2, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NoError(t, CheckPassword(password, hashedPassword2))

	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
