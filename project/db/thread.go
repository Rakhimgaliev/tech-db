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

func CreateThread(conn *pgx.ConnPool, thread *models.Thread) error {
	var user models.User
	user.Nickname = thread.Author
	err := GetUserByNickname(conn, &user)
	if err != nil {
		return err
	}
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

	nullSlug := sql.NullString{
		String: "",
		Valid:  false,
	}
	if thread.Slug == "" {
		err = transaction.QueryRow(createThread, nullSlug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &nullSlug, &thread.Created)
	} else {
		err = transaction.QueryRow(createThread, thread.Slug, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

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
