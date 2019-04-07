package db

import (
	"strings"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

const (
	createPostsRequestBeginning = `
		INSERT INTO post (forum, author, created, message, edited, parent_id, thread_id)
			VALUES 
	`

	createPostsRequestEnd = `
		RETURNING id, author, created, edited, message, parent_id, thread_id, forum
	`

	insertForumUsersStart = `
		INSERT INTO forum (nickname, forum)
			VALUES
	`

	insertForumUsersEnd = `
		ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`
)

func CreatePosts(conn *pgx.ConnPool, posts models.Posts) error {

	return nil
}

func generateCreatePostsRequest(posts models.Posts) (string, error) {
	result := strings.Builder{}

	result.WriteString(createPostsRequestBeginning)

	for i := 0; i <= len(posts); i++ {

	}

	return result.String(), nil
}
