package db

import (
	"database/sql"
	"log"
	"os"

	// "simplebank/api"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

var dbConn *sql.DB

// func newTestServer(t *testing.T, store Store) *api.Server {
// 	config := util.Config{
// 		TokenSymmetricKey: util.RandomString(32),
// 		AccessTokenDuration: time.Minute,
// 	}
// 	server, err := api.NewServer(config, store )
// 	require.NoError(t, err)

// 	return server
// 	// server, err :=
// }

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config")
	}

	dbConn, err = sql.Open(config.DBDRiver, config.DBSource)

	if err != nil {
		log.Fatal("cannot create db:", err)
	}

	testQueries = New(dbConn)

	os.Exit(m.Run())

}
