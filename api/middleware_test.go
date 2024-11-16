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

func addAuth(t *testing.T, req *http.Request, authenticationTypeBearer string, username string, duration time.Duration, tokenMaker token.Maker) {
	token, err := tokenMaker.CreateToken(username, duration)

	require.NoError(t, err)
	authToken := fmt.Sprintf("%s %s", authenticationTypeBearer, token)

	req.Header.Set(authorizationHeader, authToken)
}

func TestAuthMiddleware(t *testing.T) {
	testcases := []struct {
		name          string
		checkresponse func(*testing.T, httptest.ResponseRecorder)
		setupAuth     func(*testing.T, *http.Request, token.Maker)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", "user", time.Minute, tokenMaker)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name: "NoAuth",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				// addAuth(t, req, "user", time.Minute, tokenMaker)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},

		{
			name: "ExpiredAuth",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", "user", -time.Minute, tokenMaker)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},

		{
			name: "EmptyBearerType",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "", "user", time.Minute, tokenMaker)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},

		{
			name: "EmptyBearerType",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "invalidBearer", "user", time.Minute, tokenMaker)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
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
