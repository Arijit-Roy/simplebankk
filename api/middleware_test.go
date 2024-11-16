package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simplebank/token"
	"simplebank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(t *testing.T, request *http.Request, authType string, tokenMaker token.Maker, username string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(username, duration)

	require.NoError(t, err)
	authHeader := fmt.Sprintf("%s %s", authType, token)

	request.Header.Set(authorizationHeader, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkresponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, authTypeBearer, tokenMaker, "user", time.Minute)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "NoAuth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// addAuth(t, request, authTypeBearer, tokenMaker, "user", time.Minute)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name: "UnsupprtedAuth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, "unsupprtef", tokenMaker, "user", time.Minute)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "InvalidAuth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, "", tokenMaker, "user", time.Minute)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "ExpiredAuth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, authTypeBearer, tokenMaker, "user", -time.Minute)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config := util.Config{
				TokenSymmetricKey:   util.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			server, err := NewServer(config, nil)
			require.NoError(t, err)
			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)

			require.NoError(t, err)
			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			tc.checkresponse(t, *recorder)
		})
	}
}
