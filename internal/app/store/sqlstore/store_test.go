package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=192.168.1.21 user=tester password=tester dbname=tester sslmode=disable"
	}

	os.Exit(m.Run())
}
