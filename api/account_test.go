package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user := util.RandomString(10)
	account := randomAccount(user)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(*testing.T, *http.Request, token.Maker)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", user, time.Minute, tokenMaker)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchBodytoAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NOTFOUND",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", user, time.Minute, tokenMaker)
			},
		},

		{
			name:      "INVALID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", user, time.Minute, tokenMaker)
			},
		},

		{
			name:      "Unauthorized",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", "aaaa", time.Minute, tokenMaker)
			},
		},

		{
			name:      "Unauthorized1",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, "bearer", "", time.Minute, tokenMaker)
			},
		},

		{
			name:      "Unauthorized2",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				// addAuth(t, req, "bearer", "", time.Minute, tokenMaker)
			},
		},
	}
	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		config := util.Config{
			TokenSymmetricKey:   util.RandomString(32),
			AccessTokenDuration: time.Minute,
		}

		server, err := NewServer(config, store)

		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/account/%d", tc.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		tc.setupAuth(t, request, server.tokenMaker)
		// addAuth(t, request, "bearer", user, time.Minute, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, *recorder)
	}
}

func matchBodytoAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var req db.Account

	err = json.Unmarshal(data, &req)
	require.NoError(t, err)

	require.Equal(t, req, account)
}

func randomAccount(user string) db.Account {
	if len(user) == 0 {
		user = util.RandomOwner()
	}
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    user,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
