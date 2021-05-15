package main

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Env struct {
	dbPool *pgxpool.Pool
}

type response struct {
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	Random    int       `json:"random"`
}
