package db

import (
	"errors"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"

	"github.com/jackc/pgx"
)

var (
	ErrorUserNotFound       = errors.New("User not found")
	ErrorForumAlreadyExists = errors.New("Forum already exists")
	ErrorForumNotFound      = errors.New("Forum not found")
)

const (
	createForum = `
		INSERT INTO forum (title, userNickname, slug)
			VALUES (
			$1,
			(SELECT u.nickname FROM "user" u WHERE u.nickname = $2),
			$3)
			RETURNING title, userNickname, slug, postCount, threadCount
		`

	getForumBySlug = `
		SELECT title, userNickname, slug, postCount, threadCount
			FROM forum WHERE slug = $1`
)

const (
	PgxErrorUniqueViolation      = "23505"
	PgxErrorCodeNotNullViolation = "23502"
	PgxErrorForeignKeyViolation  = "23503"
)

func CreateForum(conn *pgx.ConnPool, forum *models.Forum) error {
	err := conn.QueryRow(createForum, forum.Title, forum.User, forum.Slug).
		Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	log.Println("CREATEFORUM: ", err)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case PgxErrorUniqueViolation:
				return ErrorForumAlreadyExists
			case PgxErrorCodeNotNullViolation:
				return ErrorUserNotFound
			}
		}
		return err
	}

	return nil
}

func GetForumBySlug(conn *pgx.ConnPool, forum *models.Forum) error {
	err := conn.QueryRow(getForumBySlug, (*forum).Slug).
		Scan(&(*forum).Title, &(*forum).User, &(*forum).Slug, &(*forum).Posts, &(*forum).Threads)
	log.Println(err)
	if err == pgx.ErrNoRows {
		return ErrorForumNotFound
	}
	return nil
}

func ForumExistsBySlug(conn *pgx.ConnPool, slug string) bool {
	err := conn.QueryRow(`SELECT FROM forum WHERE slug = $1`, slug).Scan()
	if err != nil {
		return false
	}
	return true
}
