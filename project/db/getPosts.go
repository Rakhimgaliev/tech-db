package db

import (
	"log"
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

func GetPosts(conn *pgx.ConnPool, slug_or_id string,
	limit int, desc bool, since int,
	sort string, posts *models.Posts) error {

	thread := models.Thread{}
	id, err := strconv.Atoi(slug_or_id)
	if err == nil {
		thread.Id = int32(id)
		err = GetThreadById(conn, &thread)
		if err != nil {
			return ErrorThreadNotFound
		}
	} else {
		thread.Slug = slug_or_id
		err = GetThreadBySlug(conn, &thread)
		if err != nil {
			return ErrorThreadNotFound
		}
	}

	rows, err := GenerateGetPostsQuery(conn, thread.Id, limit, desc, since, sort)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		err := scanPostRows(rows, post)

		if err != nil {
			return err
		}
		*posts = append(*posts, post)
	}

	return nil
}

const (
	getPostsFlatLimitById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1
			ORDER BY p.created, p.id
			LIMIT $2
	`

	getPostsFlatLimitDescById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1
			ORDER BY p.created DESC, p.id DESC
			LIMIT $2
	`

	getPostsFlatLimitSinceById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.id > $2
			ORDER BY p.created, p.id
			LIMIT $3
	`

	getPostsFlatLimitSinceDescById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.id < $2
			ORDER BY p.created DESC, p.id DESC
			LIMIT $3
	`

	getPostsTreeLimitById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1
			ORDER BY p.children
			LIMIT $2
	`

	getPostsTreeLimitDescById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1
			ORDER BY children DESC
			LIMIT $2
	`

	getPostsTreeLimitSinceById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and (p.children > (SELECT p2.children from post p2 where p2.id = $2))
			ORDER BY p.children
			LIMIT $3
	`

	getPostsTreeLimitSinceDescById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and (p.children < (SELECT p2.children from post p2 where p2.id = $2))
			ORDER BY p.children DESC
			LIMIT $3
	`

	getPostsParentTreeLimitById = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.children[1] IN (
				SELECT p2.children[1]
				FROM post p2
				WHERE p2.thread = $2 AND p2.parent IS NULL
				ORDER BY p2.children
				LIMIT $3
			)
			ORDER BY children
	`

	selectPostsParentTreeLimitDescByID = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.children[1] IN (
				SELECT p2.children[1]
				FROM post p2
				WHERE p2.parent IS NULL and p2.thread = $2
				ORDER BY p2.children DESC
				LIMIT $3
			)
			ORDER BY p.children[1] DESC, p.children[2:]
	`

	selectPostsParentTreeLimitSinceByID = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.children[1] IN (
				SELECT p2.children[1]
				FROM post p2
				WHERE p2.thread = $2 AND p2.parent IS NULL and p2.children[1] > (SELECT p3.children[1] from post p3 where p3.id = $3)
				ORDER BY p2.children
				LIMIT $4
			)
			ORDER BY p.children
	`

	selectPostsParentTreeLimitSinceDescByID = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p
			WHERE p.thread = $1 and p.children[1] IN (
				SELECT p2.children[1]
				FROM post p2
				WHERE p2.thread = $2 AND p2.parent IS NULL and p2.children[1] < (SELECT p3.children[1] from post p3 where p3.id = $3)
				ORDER BY p2.children DESC
				LIMIT $4
			)
			ORDER BY p.children[1] DESC, p.children[2:]
	`
)

func GenerateGetPostsQuery(conn *pgx.ConnPool,
	id int32, limit int, desc bool,
	since int, sort string) (*pgx.Rows, error) {

	var rows *pgx.Rows
	var err error

	log.Println("id:", id, "limit:", limit, desc, since, sort)

	switch sort {
	case "":
		fallthrough
	case "flat":
		if since == 0 {
			if !desc {
				rows, err = conn.Query(getPostsFlatLimitById, id, limit)
			} else {
				rows, err = conn.Query(getPostsFlatLimitDescById, id, limit)
			}
		} else {
			if !desc {
				rows, err = conn.Query(getPostsFlatLimitSinceById, id, since, limit)
			} else {
				rows, err = conn.Query(getPostsFlatLimitSinceDescById, id, since, limit)
			}
		}
	case "tree":
		if since == 0 {
			if !desc {
				rows, err = conn.Query(getPostsTreeLimitById, id, limit)
			} else {
				rows, err = conn.Query(getPostsTreeLimitDescById, id, limit)
			}
		} else {
			if !desc {
				rows, err = conn.Query(getPostsTreeLimitSinceById, id, since, limit)
			} else {
				rows, err = conn.Query(getPostsTreeLimitSinceDescById, id, since, limit)
			}
		}
	case "parent_tree":
		if since == 0 {
			if !desc {
				rows, err = conn.Query(getPostsParentTreeLimitById, id, id,
					limit)
			} else {
				rows, err = conn.Query(selectPostsParentTreeLimitDescByID, id, id,
					limit)
			}
		} else {
			if !desc {
				rows, err = conn.Query(selectPostsParentTreeLimitSinceByID, id, id,
					since, limit)
			} else {
				rows, err = conn.Query(selectPostsParentTreeLimitSinceDescByID, id, id,
					since, limit)
			}
		}
	}
	return rows, err
}
