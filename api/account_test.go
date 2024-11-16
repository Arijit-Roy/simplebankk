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
	"simplebank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
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
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// matchBodytoAccount(t, recorder.Body, account)
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
		server.router.ServeHTTP(recorder, request)
		require.NoError(t, err)
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

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
