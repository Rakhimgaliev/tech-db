package db

import (
	"database/sql"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

const (
	createThread = `
		INSERT INTO thread (slug, title, userNickname, message, created, forum)
			VALUES (
				$1, $2,
				(SELECT u.nickname FROM "user" u WHERE u.nickname = $3),
				$4, $5,
				(SELECT f.slug FROM forum f WHERE f.slug = $6)
			)
			RETURNING id, title, userNickname, forum, message, votes, slug, created
	`
)

func CreateThread(conn *pgx.ConnPool, thread *models.Thread) error {
	var err error
	if thread.Slug == "" {
		nullSlug := sql.NullString{
			String: "",
			Valid:  false,
		}
		err = conn.QueryRow(createThread, nullSlug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &nullSlug, &thread.Created)
	} else {
		err = conn.QueryRow(createThread, thread.Slug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

	log.Print("ASDASDA", err)
	if err != nil {
		return err
	}

	return nil
}
