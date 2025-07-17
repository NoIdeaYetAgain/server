package db

import (
	"clove/util"
	"context"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	defer func() {

		if err := conn.Close(context.Background()); err != nil {
			log.Fatal("failed to close connection", err)
		}
	}()

	testQueries = New(conn)

	os.Exit(m.Run())
}
