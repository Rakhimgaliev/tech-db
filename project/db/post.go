package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

var (
	ErrorThreadNotFound       = errors.New("Forum already exists")
	ErrorPostCreateConflict   = errors.New("Post create conflict")
	ErrorPostCreateBadRequest = errors.New("Post Create Bad Request")
)

const (
	getForumSlugByThreadId = `
		SELECT forum FROM thread WHERE id = $1
	`

	getForumSlugAndThreadIdByThreadSlug = `
		SELECT forum, id from thread WHERE slug = $1
	`

	createPostsQueryStart = `
		INSERT INTO post (forum, userNickname, created, message, edited, parent, thread)
		VALUES 
	`

	createPostsQueryEnd = `
		RETURNING id, userNickname, created, edited, message, parent, thread, forum
	`

	createForumUsersQueryStart = `
		INSERT INTO forum_user (nickname, forum)
		VALUES
	`

	createForumUsersQueryEnd = `
		ON CONFLICT ON CONSTRAINT uniqueForumUser DO NOTHING
	`

	createWithCheckParentID = `(
		SELECT (
			CASE WHEN 
			EXISTS(SELECT 1 from post p where p.id=%v and p.thread=%v)
			THEN %v ELSE -1 END)
		)
	`

	updateForumPostCountByThreadId = `
		UPDATE forum f SET postCount = postCount + $1
		FROM thread t
		WHERE t.forum = f.slug AND t.id = $2
	`
)

func CreatePosts(conn *pgx.ConnPool, threadIdOrSlag string, posts *models.Posts) error {
	if len(*posts) == 0 {
		return ErrorPostCreateBadRequest
	}

	forumSlug, threadId, err := getForumSlugAndThreadIdByThreadSlugOrId(conn, threadIdOrSlag)
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

	err = createPosts(transaction, threadId, posts, forumSlug)

	if err != nil {
		if txErr := transaction.Rollback(); txErr != nil {
			return txErr
		}

		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case PgxErrorForeignKeyViolation:
				if pqError.ConstraintName == "post_parent_id_fkey" {
					return ErrorPostCreateConflict
				}
				if pqError.ConstraintName == "post_author_fkey" {
					return ErrorUserNotFound
				}
			}
		}
		return err
	}

	if err = transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func createPosts(transaction *pgx.Tx, threadId int32, posts *models.Posts, forumSlug string) error {
	postsArgs := make([]interface{}, 0)
	forumUserArgs := make([]interface{}, 0)

	createPostsQuery, createForumUserQuery := generatePostsQuery(forumSlug, threadId, posts, &postsArgs, &forumUserArgs)

	*posts = nil
	rows, err := transaction.Query(*createPostsQuery, postsArgs...)
	if err != nil {
		return err
	}
	for rows.Next() {
		post := &models.Post{}
		err := scanPostRows(rows, post)
		if err != nil {
			rows.Close()
			return err
		}

		*posts = append(*posts, post)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}

	rows.Close()

	_, err = transaction.Exec(updateForumPostCountByThreadId, len(*posts), threadId)
	if err != nil {
		return err
	}

	_, err = transaction.Exec(*createForumUserQuery, forumUserArgs...)
	return err
}

func scanPostRows(r *pgx.Rows, post *models.Post) error {
	parent := sql.NullInt64{}
	err := r.Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}
	return err
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
	values := "("
	valuesArr := []string{}
	placeholder := placeholderStart

	valuesArr = append(valuesArr, fmt.Sprintf(`'%v'`, forumSlug))

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if author == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, author)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if created.IsZero() {
		*valuesArgs = append(*valuesArgs, "now()")
	} else {
		*valuesArgs = append(*valuesArgs, created)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if message == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, message)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	*valuesArgs = append(*valuesArgs, isEdited)

	if parent == 0 {
		valuesArr = append(valuesArr, fmt.Sprint("(NULL)"))
	} else {
		valuesArr = append(valuesArr, fmt.Sprintf(createWithCheckParentID, parent, threadId, parent))
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	if postThreadId == 0 {
		*valuesArgs = append(*valuesArgs, threadId)
	} else {
		*valuesArgs = append(*valuesArgs, postThreadId)
	}

	values += strings.Join(valuesArr, ", ")
	values += ")"

	return values
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
	err := conn.QueryRow(getForumSlugAndThreadIdByThreadSlug, threadIdOrSlug).Scan(&forumSlug, &threadId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forumSlug, threadId, ErrorThreadNotFound
		}
		return forumSlug, threadId, err
	}

	return forumSlug, threadId, nil
}
