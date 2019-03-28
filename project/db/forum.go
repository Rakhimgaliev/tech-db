package db

import (
	"errors"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"

	"github.com/jackc/pgx"
)

var (
	ErrorUserNotFound      = errors.New("User not found")
	ErrorForumAlreadyExist = errors.New("Forum already exists")
	ErrorForumNotFound     = errors.New("Forum not found")
)

const (
	createForum = `
		INSERT INTO forum (user_nick, slug, title) 
			VALUES 
			((SELECT u.nickname FROM "user" u WHERE u.nickname = $1),
			$2, $3)
			RETURNING userNickname, slug, title, postCount, threadCount
		`

	getForum = `SELECT FROM forum WHERE slig = $1`
)

const (
	PgxErrorUniqueViolation      = "23505"
	PgxErrorForeignKeyViolation  = "23503"
	PgxErrorCodeNotNullViolation = "23502"
)

func CreateForum(conn *pgx.ConnPool, forum *models.Forum) error {
	err := conn.QueryRow(createForum, forum.User, forum.Slug, forum.Title).Scan(forum)
	log.Println(err)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case PgxErrorUniqueViolation:
				return ErrorUserNotFound
			case PgxErrorCodeNotNullViolation:
				return ErrorForumAlreadyExist
			}
		}
		return err
	}

	return nil
}

func GetForumBySlug(conn *pgx.ConnPool, forum *models.Forum) error {
	err := conn.QueryRow(getForum, forum.Slug).Scan(forum)
	if err == pgx.ErrNoRows {
		return ErrorForumNotFound
	}
	return nil
}

// const (
// 	checkForumExist = `SELECT FROM forum WHERE slug = $1`
// 	checkUserExist  = `SELECT FROM "user" WHERE nickname = $1`
// )

// func CheckForumExist(conn *pgx.ConnPool, forumSlug string) (bool, error) {
// 	err := conn.QueryRow(checkForumExist, forumSlug).Scan()
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }

// func CheckUserExist(conn *pgx.ConnPool, userNickname string) (bool, error) {
// 	err := conn.QueryRow(checkUserExist, userNickname).Scan()
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }
