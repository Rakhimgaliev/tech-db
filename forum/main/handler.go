package main

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db-forum/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

func Status(c *gin.Context, conn *pgx.ConnPool) {
	var status models.Status
	conn.QueryRow(statusQuery).Scan(
		&status.Forum, &status.Thread,
		&status.Post, &status.User)
	answer, _ := json.Marshal(status)
	c.JSON(200, string(answer))
}
