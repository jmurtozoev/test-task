package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	dbConn *sqlx.DB
	once   sync.Once
)

func Connect() *sqlx.DB {
	once.Do(func() {
		db, err := sqlx.Connect("postgres", getDSN())
		if err != nil {
			log.Fatal(err)
		}
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(9)

		dbConn = db
	})

	return dbConn
}

// TODO move to config
func getDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 	// cfg.PostgresHost,
		 	5432,			// cfg.PostgresPort,
			"postgres",		// cfg.PostgresUser,
			"1234",			// cfg.PostgresPassword,
		    "test", 		// cfg.PostgresDatabase,
		)
}
