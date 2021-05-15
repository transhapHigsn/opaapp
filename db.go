package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

func initDb() *pgxpool.Pool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	pid := os.Getpid()

	user := os.Getenv("OPAAPP_DB_USER")
	password := os.Getenv("OPAAPP_DB_PASSWORD")
	host := os.Getenv("OPAAPP_DB_HOST")
	dbport := os.Getenv("OPAAPP_DB_PORT")
	dbname := os.Getenv("OPAAPP_DB")

	port, err := strconv.ParseInt(dbport, 10, 64)
	if err != nil {
		log.Printf("pid=%d func=initDb level=error msg=Failed reading DB configuration.", pid)
		panic(err)
	}

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
