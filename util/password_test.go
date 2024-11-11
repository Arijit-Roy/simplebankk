package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	passwd := RandomString(10)
	hashedPassword, err := HashPassword(passwd)
	require.NoError(t, err)

	err = CheckPassword(passwd, hashedPassword)
	require.NoError(t, err)

	hashedPassword1, err := HashPassword(passwd)
	require.NoError(t, err)

	require.NotEqual(t, hashedPassword, hashedPassword1)

}
