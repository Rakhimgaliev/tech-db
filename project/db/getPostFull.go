package db

import (
	"database/sql"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

func GetPostFull(conn *pgx.ConnPool, related []string, postFull *models.PostFull) error {
	withUser, withThread, withForum := false, false, false
	for _, rel := range related {
		switch rel {
		case "user":
			postFull.Author = &models.User{}
			withUser = true
		case "thread":
			postFull.Thread = &models.Thread{}
			withThread = true
		case "forum":
			postFull.Forum = &models.Forum{}
			withForum = true
		}
	}

	log.Println("withUser, withThread, withForum ", withUser, withThread, withForum)

	var err error
	if !withUser && !withThread && !withForum {
		err = getPost(conn, postFull.Post)
	} else if withUser && !withThread && !withForum {
		err = getPostWithUser(conn, postFull)
	} else if !withUser && withThread && !withForum {
		err = getPostWithThread(conn, postFull)
	} else if !withUser && !withThread && withForum {
		err = getPostWithForum(conn, postFull)
	} else if withUser && withThread && !withForum {
		err = getPostWithUserThread(conn, postFull)
	}

	log.Println("----------", err)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrorPostNotFound
		}
		return err
	}
	return nil
}

const (
	getPostQuery = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum
			FROM post p 
			WHERE p.id = $1
	`

	getPostWithUserQuery = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum,
			u.nickname, u.fullname, u.about, u.email
			FROM post p 
			JOIN "user" u ON p.userNickname = u.nickname
			WHERE p.id = $1
	`

	getPostWithThreadQuery = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum,
			t.id, t.slug, t.userNickname, t.created, t.forum, t.title, t.message, t.votes
			FROM post p 
			JOIN thread t ON p.thread = t.id
			WHERE p.id = $1
	`

	getPostWithForumQuery = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum,
			f.userNickname, f.slug, f.title, f.threadCount, f.postCount
			FROM post p 
			JOIN thread t ON p.thread = t.id
			JOIN forum f ON p.forum = f.slug
			WHERE p.id = $1
	`

	getPostWithUserThreadQuery = `
		SELECT p.id, p.userNickname, p.created, p.edited, p.message, p.parent, p.thread, p.forum,
			t.id, t.slug, t.userNickname, t.created, t.forum, t.title, t.message, t.votes,
			u.nickname, u.fullname, u.about, u.email
			FROM post p 
			JOIN thread t ON p.thread = t.id
			JOIN "user" u ON p.userNickname = u.nickname
			WHERE p.id = $1
	`
)

func getPost(conn *pgx.ConnPool, post *models.Post) error {
	parent := sql.NullInt64{}
	err := conn.QueryRow(getPostQuery, post.Id).Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)
	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}
	return err
}

func getPostWithUser(db *pgx.ConnPool, postFull *models.PostFull) error {
	parent := sql.NullInt64{}
	err := db.QueryRow(getPostWithUserQuery, postFull.Post.Id).Scan(
		&postFull.Post.Id,
		&postFull.Post.Author,
		&postFull.Post.Created,
		&postFull.Post.IsEdited,
		&postFull.Post.Message,
		&parent,
		&postFull.Post.Thread,
		&postFull.Post.Forum,
		&postFull.Author.Nickname,
		&postFull.Author.Fullname,
		&postFull.Author.About,
		&postFull.Author.Email,
	)
	if err != nil {
		return err
	}

	if parent.Valid {
		postFull.Post.Parent = parent.Int64
	} else {
		postFull.Post.Parent = 0
	}
	return nil
}

func getPostWithThread(db *pgx.ConnPool, postFull *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(getPostWithThreadQuery, postFull.Post.Id).Scan(
		&postFull.Post.Id,
		&postFull.Post.Author,
		&postFull.Post.Created,
		&postFull.Post.IsEdited,
		&postFull.Post.Message,
		&parent,
		&postFull.Post.Thread,
		&postFull.Post.Forum,
		&postFull.Thread.Id,
		&slugThread,
		&postFull.Thread.Author,
		&postFull.Thread.Created,
		&postFull.Thread.Forum,
		&postFull.Thread.Title,
		&postFull.Thread.Message,
		&postFull.Thread.Votes,
	)
	if err != nil {
		return err
	}

	if parent.Valid {
		postFull.Post.Parent = parent.Int64
	} else {
		postFull.Post.Parent = 0
	}
	if slugThread.Valid {
		postFull.Thread.Slug = slugThread.String
	} else {
		postFull.Thread.Slug = ""
	}

	return nil
}

func getPostWithForum(db *pgx.ConnPool, postFull *models.PostFull) error {
	parent := sql.NullInt64{}
	err := db.QueryRow(getPostWithForumQuery, postFull.Post.Id).Scan(
		&postFull.Post.Id,
		&postFull.Post.Author,
		&postFull.Post.Created,
		&postFull.Post.IsEdited,
		&postFull.Post.Message,
		&parent,
		&postFull.Post.Thread,
		&postFull.Post.Forum,
		&postFull.Forum.User,
		&postFull.Forum.Slug,
		&postFull.Forum.Title,
		&postFull.Forum.Threads,
		&postFull.Forum.Posts,
	)
	if err != nil {
		return err
	}

	if parent.Valid {
		postFull.Post.Parent = parent.Int64
	} else {
		postFull.Post.Parent = 0
	}
	return nil
}

func getPostWithUserThread(db *pgx.ConnPool, pf *models.PostFull) error {
	parent := sql.NullInt64{}
	slugThread := sql.NullString{}
	err := db.QueryRow(getPostWithUserThreadQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Thread.Id,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	if err != nil {
		return err
	}

	if parent.Valid {
		pf.Post.Parent = parent.Int64
	} else {
		pf.Post.Parent = 0
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}
