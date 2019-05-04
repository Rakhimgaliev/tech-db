package db

import (
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

const (
	createThreadVote = `
		INSERT INTO vote (nickname, voice, threadId)
			VALUES ($1, $2, $3)
			ON CONFLICT ON CONSTRAINT uniqueVote 
			DO UPDATE SET voice = EXCLUDED.voice;
	`

	getThreadIdBySlug = `
		SELECT id from thread WHERE slug = $1
	`

	updateThreadVotesCountQuery = `
		UPDATE thread t SET votes = (
			SELECT SUM(case when v.voice = true then 1 else -1 end)
				FROM vote v 
				WHERE v.threadId=$1
		)
			WHERE t.id=$2
				RETURNING id, title, userNickname, forum, message, votes, slug, created
	`
)

func CreateThreadVote(conn *pgx.ConnPool, threadSlugOrId string, thread *models.Thread, vote *models.Vote) error {
	if threadId, err := strconv.Atoi(threadSlugOrId); err != nil {
		threadID, err := GetThreadIdBySlug(conn, threadSlugOrId)

		if err != nil {
			if err == pgx.ErrNoRows {
				return ErrorThreadNotFound
			}
			return err
		}
		thread.Id = int32(threadID)
	} else {
		thread.Id = int32(threadId)
	}

	var voteBool = false
	if vote.Voice == 1 {
		voteBool = true
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

	_, err = transaction.Exec(createThreadVote, vote.Nickname, voteBool, thread.Id)
	if err != nil {
		if txErr := transaction.Rollback(); txErr != nil {
			return txErr
		}
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case PgxErrorForeignKeyViolation:
				return ErrorThreadNotFound
			}
		}
		return err
	}

	err = updateThreadVotesCount(transaction, thread)
	if err != nil {
		if txErr := transaction.Rollback(); txErr != nil {
			return txErr
		}
		return err
	}

	if commitErr := transaction.Commit(); commitErr != nil {
		return commitErr
	}

	return nil
}

func GetThreadIdBySlug(conn *pgx.ConnPool, threadSlug string) (int, error) {
	var threadId int
	err := conn.QueryRow(getThreadIdBySlug, threadSlug).Scan(&threadId)
	return threadId, err
}

func updateThreadVotesCount(tx *pgx.Tx, thread *models.Thread) error {
	return scanThread(tx.QueryRow(updateThreadVotesCountQuery, thread.Id, thread.Id), thread)
}
