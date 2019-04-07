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

	getThreads = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE forum = $1
			ORDER BY created
			LIMIT $2
	`

	getThreadsSince = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE forum = $1 AND created <= $2
			ORDER BY created
			LIMIT $3
	`

	getThreadsDesc = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE forum = $1
			ORDER BY created DESC
			LIMIT $2	
	`

	getThreadsDescSince = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE forum = $1 AND created <= $2
			ORDER BY created DESC
			LIMIT $3
	`

	maxLimit = 9223372036854775807
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

	if err != nil {
		return err
	}

	return nil
}

func GetThreads(conn *pgx.ConnPool, slug string, limit int, since string, desc bool, threads *models.Threads) error {
	if !ForumExistsBySlug(conn, slug) {
		return ErrorForumNotFound
	}

	if limit == 0 {
		limit = maxLimit
	}

	var rows *pgx.Rows
	var err error
	if desc {
		if len(since) <= 1 {
			rows, err = conn.Query(getThreadsDesc, slug, limit)
		} else {
			rows, err = conn.Query(getThreadsDescSince, slug, since, limit)
		}
	} else {
		if len(since) <= 1 {
			rows, err = conn.Query(getThreads, slug, limit)
		} else {
			rows, err = conn.Query(getThreadsSince, slug, since, limit)
		}
	}

	if err != nil {
		log.Println("rows scan error: ", err)
		return err
	}

	defer rows.Close()
	for rows.Next() {
		thread := &models.Thread{}
		nullableSlug := sql.NullString{}
		err := rows.Scan(&thread.Id, &nullableSlug, &thread.Author, &thread.Created, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			return err
		}
		if nullableSlug.Valid {
			thread.Slug = nullableSlug.String
		}
		*threads = append(*threads, thread)
	}

	return nil
}
