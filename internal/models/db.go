package models

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error{
	var err error
	connStr := "postgres://postgres:postgres@db:5432/pvz_db" 

	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return err
	}

	DB = pool
	return nil
}