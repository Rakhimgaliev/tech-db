package db

import (
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

func GetPosts(conn *pgx.ConnPool, slug_or_id string,
	limit string, desc bool, since int,
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

	rows, err := GenerateGetPostsQuery(conn, id, limit, desc, since, sort)
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

const getPostsbeginning = `
	SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path
		LIMIT $3
	)
	ORDER BY path
`

func GenerateGetPostsQuery(conn *pgx.ConnPool,
	id int, limit string, desc bool,
	since int, sort string) (*pgx.Rows, error) {

	var rows *pgx.Rows
	var err error

	switch sort {
	case "":
		fallthrough
	case "flat":

	case "tree":

	case "parent_tree":
	}
	return rows, err
}
