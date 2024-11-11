package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

var dbConn *sql.DB

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
