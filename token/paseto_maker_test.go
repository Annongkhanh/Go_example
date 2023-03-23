package token

import (
	"testing"
	"time"

	"github.com/Annongkhanh/Go_example/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := util.RandomString(32)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)

	require.NoError(t, err)

	username := util.RandomString(8)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt.Add(time.Minute), payload.ExpiredAt, time.Second)
	require.NotZero(t, payload.ID)

}

func TestExpiredPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	username := util.RandomString(8)
	duration := time.Minute

	token, payload, err := pasetoMaker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}
