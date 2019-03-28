package handler

import (
	"log"

	"github.com/jackc/pgx"
)

type handler struct {
	conn *pgx.ConnPool
}

func NewConnPool(config *pgx.ConnConfig) *handler {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     *config,
		MaxConnections: 3,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &handler{
		conn: connPool,
	}
}
