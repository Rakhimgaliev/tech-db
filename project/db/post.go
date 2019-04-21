package db

import (
	"errors"
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

var (
	ErrorThreadNotFound = errors.New("Forum already exists")
)

const (
	getForumSlugByThreadId = `
		SELECT forum FROM thread WHERE id = $1
	`
	getForumSlugAndThreadIdByThreadSlug = `
		Select forum, id from thread WHERE slug = $1
	`
)

func CreatePosts(conn *pgx.ConnPool, threadIdOrSlag string, posts *models.Posts) error {
	forumSlug, threadId, err := GetForumSlugAndThreadIdByThreadSlugOrId(conn, threadIdOrSlag)
	if err != nil {
		return err
	}

	transaction, err := conn.Begin()
	if err != nil {
		return err
	}

	err = insertPosts(transaction, threadId, posts, forumSlug)
	if err != nil {

	}

	return nil
}

func GetForumSlugAndThreadIdByThreadSlugOrId(conn *pgx.ConnPool, threadIdOrSlug string) (string, int, error) {
	threadId := -1
	forumSlug := ""
	if threadId, err := strconv.Atoi(threadIdOrSlug); err == nil {
		err := conn.QueryRow(getForumSlugByThreadId, threadId).Scan(&forumSlug)
		if err != nil {
			if err == pgx.ErrNoRows {
				return forumSlug, threadId, ErrorForumAlreadyExists
			}
			return forumSlug, threadId, err
		}
		return forumSlug, threadId, nil
	}
	err := conn.QueryRow(getForumSlugAndThreadIdByThreadSlug, threadIdOrSlug).Scan(&threadId, &forumSlug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forumSlug, threadId, ErrorThreadNotFound
		}
		return forumSlug, threadId, err
	}

	return forumSlug, threadId, nil
}
