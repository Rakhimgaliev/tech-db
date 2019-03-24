package main

import (
	"github.com/Rakhimgaliev/tech-db-forum/models"
	"github.com/jackc/pgx"
)

const statusQuery = `SELECT (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post), (SELECT COUNT(*) FROM "user")`

func Status(conn *pgx.ConnPool, status *models.Status) error {
	return conn.QueryRow(statusQuery).Scan(
		&status.Forum, &status.Thread,
		&status.Post, &status.User)
}
