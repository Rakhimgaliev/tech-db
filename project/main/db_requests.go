package main

const (
	createForum = `
		INSERT INTO FORUM (userNickname, slug, title)
		VALUE(
			(SELECT u.nickname FROM "user" u WHRE u.nickname = $1),
			$2,
			$3)
		RETUURNING userNickname, slug, title, threadCount, postCount
		`

	checkForumExist = `
		SELECT FROM forum where slug = $1
		`

	forumDetails = `
		SELECT user_nick, slug, title, threadCount, postCount 
			FROM forum WHERE slug = $1
		`

	createThread = `
		INSERT INTO thread (slug, userNickname, created, forumSlug, title, message)
		VALUES(
			$1,
			(SELECT u.nickname FROM "user" u WHERE u.nickname = $2),
			$3,
			(SELECT f.slug FROM forum f WHERE f.slug = $4),
			$5,
			$6)
		RETURNING id, slug, userNickname, created, forumSlug, title, message, votes
		`

	getUsersByForum = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname
		`

	getUsersByForumLimit = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname
			LIMIT $2
		`

	getUsersByForumSince = `
		SELECT u.nickname, u.fullname, u.about,u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname
		`

	getUsersByForumDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname DESC
		`

	getUsersByForumLimitSince = `
		SELECT u.nickname, u.fullname, u.about,u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname
			LIMIT $3
		`

	getUsersByForumLimitDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname DESC
			LIMIT $2
		`

	getUsersByForumSinceDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname DESC
		`

	getUsersByForumLimitSinceDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname DESC
			LIMIT $3
		`

	getThreadsByForum = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname
		`

	getThreadsByForumLimit = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname
			LIMIT $2
		`

	getThreadsByForumSince = `
		SELECT u.nickname, u.fullname, u.about,u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname
		`

	getThreadsByForumDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname DESC
		`

	getThreadsByForumLimitSince = `
		SELECT u.nickname, u.fullname, u.about,u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname
			LIMIT $3
		`

	getThreadsByForumLimitDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1
			ORDER BY u.nickname DESC
			LIMIT $2
		`

	getThreadsByForumSinceDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname DESC
		`

	getThreadsByForumLimitSinceDesc = `
		SELECT u.nickname, u.fullname. u.about, u.email
			FROM forum_user f_u
			JOIN "user" u on f_u.nickname = u.nickname
			WHERE f_u.forumSlug = $1 AND u.nickname > $2
			ORDER BY u.nickname DESC
			LIMIT $3
		`

	statusQuery = `
		SELECT
			(SELECT COUNT(*) FROM forum),
			(SELECT COUNT(*) FROM thread),
			(SELECT COUNT(*) FROM post),
			(SELECT COUNT(*) FROM "user")
		`
)
