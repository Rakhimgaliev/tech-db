package db

import "errors"

var (
	ErrorUserNotFound      = errors.New("User not found")
	ErrorForumAlreadyExist = errors.New("")
)

const createForum = `
		INSERT INTO FORUM (userNickname, slug, title)
		VALUE(
			(SELECT u.nickname FROM "user" u WHRE u.nickname = $1),
			$2,
			$3)
		RETURNING userNickname, slug, title, threadCount, postCount
		`

func CreateForum() error {

	return ErrorForumAlreadyExist
}
