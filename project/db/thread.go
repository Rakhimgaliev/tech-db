package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

var (
	ErrorUniqueViolation = errors.New("Error Unique Violatation")
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

	getThreadBySlug = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE slug = $1
	`

	getThreadById = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE id = $1
	`

	getThreadsSince = `
		SELECT id, slug, userNickname, created, forum, title, message, votes
			FROM thread
			WHERE forum = $1 AND created >= $2
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

func scanThread(row *pgx.Row, thread *models.Thread) error {
	threadSlug := sql.NullString{}

	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &threadSlug, &thread.Created)
	log.Println("error on Scanning:", err)
	if err != nil {
		return err
	}
	if threadSlug.Valid {
		thread.Slug = threadSlug.String
	}
	return err
}

func CreateThread(conn *pgx.ConnPool, thread *models.Thread) error {
	transaction, err := conn.Begin()
	if err != nil {
		return err
	}

	_, err = transaction.Exec("SET LOCAL synchronous_commit TO OFF")
	if err != nil {
		if txErr := transaction.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}

	threadSlug := sql.NullString{}
	if thread.Slug == "" {
		threadSlug = sql.NullString{
			String: "",
			Valid:  false,
		}
	} else {
		threadSlug = sql.NullString{
			String: thread.Slug,
			Valid:  true,
		}
	}

	err = scanThread(transaction.QueryRow(createThread, threadSlug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum), thread)
	log.Println("HERE:", err)
	if err != nil {
		if txErr := transaction.Rollback(); txErr != nil {
			return txErr
		}
		if err, ok := err.(pgx.PgError); ok {
			switch err.Code {
			case PgxErrorUniqueViolation:
				err := GetThreadBySlug(conn, thread)
				if err == nil {
					return ErrorUniqueViolation
				}
			case PgxErrorCodeNotNullViolation:
				return ErrorUserNotFound
			}
		}
		return err
	}

	if commitErr := transaction.Commit(); commitErr != nil {
		return commitErr
	}
	return nil
}

func GetThreadBySlug(conn *pgx.ConnPool, thread *models.Thread) error {
	err := conn.QueryRow(getThreadBySlug, thread.Slug).
		Scan(&thread.Id, &thread.Slug, &thread.Author, &thread.Created, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes)
	log.Println(err)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok {
			if err.Code == PgxErrorCodeNotNullViolation {
				return ErrorThreadNotFound
			}
		}
		return err
	}
	return nil
}

func GetThreadById(conn *pgx.ConnPool, thread *models.Thread) error {
	err := conn.QueryRow(getThreadById, thread.Id).
		Scan(&thread.Id, &thread.Slug, &thread.Author, &thread.Created, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes)
	log.Println(err)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok {
			if err.Code == PgxErrorCodeNotNullViolation {
				return ErrorThreadNotFound
			}
		}
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
