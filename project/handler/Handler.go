package handler

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Rakhimgaliev/tech-db/project/db"
	"github.com/Rakhimgaliev/tech-db/project/models"
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

func (h *handler) Clear(context *gin.Context) {
	err := db.Clear(h.conn)
	if err != nil {
		return
	}

	clearJSON, _ := json.Marshal("")
	context.Data(200, "application/json", clearJSON)
	return
}

func (h *handler) Status(context *gin.Context) {
	status := models.Status{}
	err := db.Status(h.conn, &status)
	if err != nil {
		return
	}
	statusJSON, _ := json.Marshal(status)
	context.Data(200, "application/json", statusJSON)
	return
}
