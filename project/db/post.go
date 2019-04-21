package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

	createPostsQueryStart = `
		INSERT INTO post (forum, author, created, message, edited, parent, thread)
		VALUES 
	`

	createPostsQueryEnd = `
		RETURNING id, author, created, edited, message, parent, thread, forum
	`

	createForumUsersQueryStart = `
		INSERT INTO forum_user (nickname, forum)
		VALUES
	`

	createForumUsersQueryEnd = `
		ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`
)

func CreatePosts(conn *pgx.ConnPool, threadIdOrSlag string, posts *models.Posts) error {
	forumSlug, threadId, err := getForumSlugAndThreadIdByThreadSlugOrId(conn, threadIdOrSlag)
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

func insertPosts(transaction *pgx.Tx, threadId int32, posts *models.Posts, forumSlug string) error {
	if len(*posts) == 0 {
		return nil
	}

	postsArgs := make([]interface{}, 0)
	forumUserArgs := make([]interface{}, 0)
	insertPostsQuery, insertForumUserQuery := generatePostsQuery(forumSlug, threadId, posts, &postsArgs, &forumUserArgs)

	return nil
}

func generatePostsQuery(forumSlug string, threadId int32, posts *models.Posts,
	postsArgs *[]interface{}, forumUserArgs *[]interface{}) (*string, *string) {

	createValues := ""
	createUserValues := ""
	finalCreateValues := strings.Builder{}
	finalUserCreateValues := strings.Builder{}

	for idx, post := range *posts {
		createValues = formCreateValuesId(post.Author, post.Created, post.Message, threadId,
			post.IsEdited, post.Parent, post.Thread, idx*5+1, postsArgs, forumSlug)

		createUserValues = formCreateUserValues(post.Author, idx*2+1, forumUserArgs, forumSlug)

		if idx != 0 {
			finalCreateValues.WriteString(",")
			finalUserCreateValues.WriteString(",")
		}
		finalCreateValues.WriteString(createValues)
		finalUserCreateValues.WriteString(createUserValues)
	}

	createPostsQuery := strings.Builder{}
	createPostsQuery.WriteString(createPostsQueryStart)
	createPostsQuery.WriteString(finalCreateValues.String())
	createPostsQuery.WriteString(createPostsQueryEnd)

	createUsersQuery := strings.Builder{}
	createUsersQuery.WriteString(createForumUsersQueryStart)
	createUsersQuery.WriteString(finalUserCreateValues.String())
	createUsersQuery.WriteString(createForumUsersQueryEnd)

	resPosts, resUser := createPostsQuery.String(), createUsersQuery.String()
	return &resPosts, &resUser
}

func formCreateValuesId(author string, created time.Time, message string, threadId int32, isEdited bool, parent int64,
	postThreadId int32, placeholderStart int, valuesArgs *[]interface{}, forumSlug string) string {
	res := ""
	return res
}

func formCreateUserValues(author string, placeholder int, args *[]interface{}, forumSlug string) string {
	values := fmt.Sprintf("($%v, $%v)", placeholder, placeholder+1)
	*args = append(*args, author)
	*args = append(*args, forumSlug)
	return values
}

func getForumSlugAndThreadIdByThreadSlugOrId(conn *pgx.ConnPool, threadIdOrSlug string) (string, int32, error) {
	var threadId int32
	forumSlug := ""
	if threadId, err := strconv.Atoi(threadIdOrSlug); err == nil {
		err := conn.QueryRow(getForumSlugByThreadId, threadId).Scan(&forumSlug)
		if err != nil {
			if err == pgx.ErrNoRows {
				return forumSlug, int32(threadId), ErrorForumAlreadyExists
			}
			return forumSlug, int32(threadId), err
		}
		return forumSlug, int32(threadId), nil
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
