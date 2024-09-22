package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConn, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot create db:", err)
	}

	testQueries = New(dbConn)

	os.Exit(m.Run())

}
