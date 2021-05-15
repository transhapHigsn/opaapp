package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func initDb() *pgxpool.Pool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	pid := os.Getpid()

	log.Printf("pid=%d func=initDb level=info msg=Initializing Db.", pid)
	psqlconn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	dbPool, err := pgxpool.Connect(context.Background(), psqlconn)

	if err != nil {
		log.Printf("pid=%d func=initDb level=error msg=Failed initializing db.", pid)
		panic(err)
	}
	log.Printf("pid=%d func=initDb level=info msg=Initialized Db.", pid)
	return dbPool
}
