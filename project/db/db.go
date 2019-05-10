package db

import (
	"github.com/Rakhimgaliev/tech-db/project/models"
	"github.com/jackc/pgx"
)

type db struct{}

const (
	statusQuery = `
	SELECT
		(SELECT COUNT(*) FROM "user"), (SELECT COUNT(*) FROM forum),
		(SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post)
	`
	clearQuery = `TRUNCATE ONLY post, vote, thread, forum_user, forum, "user"`
)

func Clear(db *pgx.ConnPool) error {
	_, err := db.Exec(clearQuery)
	return err
}

func Status(db *pgx.ConnPool, s *models.Status) error {
	return db.QueryRow(statusQuery).Scan(&s.User, &s.Forum, &s.Thread, &s.Post)
}
