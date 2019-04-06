package db

import (
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

const (
	createThread = `
		INSERT INTO thread (slug, title, userNickname, message, created, forum_slug)
			VALUES (
				$1, $2,
				(SELECT u.nickname FROM "user" u WHERE u.nickname = $3),
				,$4, $5,
				(SELECT f.slug FROM forum f WHERE f.slug = $6)
			)
			RETURNING id, title, userNickname, forum_slug, message, votes, slug, created
	`
)

func CreateThread(conn *pgx.ConnPool, thread *models.Thread) error {
	err := conn.QueryRow(createThread, thread.Slug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	log.Print(err)
	if err != nil {
		return err
	}

	return nil
}
