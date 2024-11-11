package token

import (
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	IssuedAt := time.Now()
	expiredAt := IssuedAt.Add(duration)

	payload, err := NewPayload(username, &duration)
	require.NoError(t, err)

	err = payload.Valid()
	require.NoError(t, err)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, IssuedAt, time.Second)

	_, err = maker.CreateToken(username, duration)
	require.NoError(t, err)

}

func TestExpiredJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.Error(t, err)

	require.Empty(t, payload)

}
