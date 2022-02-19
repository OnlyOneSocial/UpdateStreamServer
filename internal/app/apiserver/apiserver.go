package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/katelinlis/AuthBackend/internal/app/store/sqlstore"
)

//Start database
func Start(config *Config) {
	fmt.Println("Hello world")

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		fmt.Print(err)
	}

	defer db.Close()

	store := sqlstore.New(db)

	srv := newServer(store, config)

	fmt.Println("Start webserver on", config.BindAddr)
	fmt.Println("Start socket on", config.BindAddrSocket)
	go srv.ListenAndServe(config.BindAddrSocket)
	err = http.ListenAndServe(config.BindAddr, srv)

}

func newDB(DatabaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
