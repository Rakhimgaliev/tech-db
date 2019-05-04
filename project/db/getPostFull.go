package db

import (
	"database/sql"

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

	var err error
	if !withForum && !withUser && !withThread {
		err = getPost(conn, postFull.Post)
	}
	// else if withForum && withUser && withThread {
	// 	err = getPostWithForumUserThread(db, pf)
	// } else if !withForum && withUser && withThread {
	// 	err = getPostWithUserThread(db, pf)
	// } else if withForum && !withUser && withThread {
	// 	err = getPostWithForumThread(db, pf)
	// } else if withForum && withUser && !withThread {
	// 	err = getPostWithForumUser(db, pf)
	// } else if !withForum && !withUser && withThread {
	// 	err = getPostWithThread(db, pf)
	// } else if !withForum && withUser && !withThread {
	// 	err = getPostWithUser(db, pf)
	// } else if withForum && !withUser && !withThread {
	// 	err = getPostWithForum(db, pf)
	// }

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
