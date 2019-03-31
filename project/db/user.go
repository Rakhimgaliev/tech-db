package db

import (
	"errors"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/jackc/pgx"
)

var (
	ErrorUserAlreadyExists = errors.New("User not found")
)

const (
	createUser = `
		INSERT INTO "user" (nickname, fullname, about, email) 
			VALUES ($1, $2, $3, $4)
			RETURNING nickname, fullname, about, email
		`

	getUserByNickname = `
		SELECT nickname, fullname, about, email
			FROM "user" WHERE nickname = $1`

	getUserByEmailOrNickname = `
		SELECT nickname, fullname, about, email
			FROM "user" WHERE email = $1 OR nickname = $2`
)

func CreateUser(conn *pgx.ConnPool, user *models.User) error {
	err := conn.QueryRow(createUser, (*user).Nickname, (*user).Fullname, (*user).About, (*user).Email).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	log.Println(err)

	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			if pqError.Code == PgxErrorUniqueViolation {
				return ErrorUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func GetUserByEmailOrNickname(conn *pgx.ConnPool, email string, nickname string) (models.Users, error) {
	rows, err := conn.Query(getUserByEmailOrNickname, email, nickname)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	users := models.Users{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

func GetUserByNickname(conn *pgx.ConnPool, user *models.User) error {
	err := conn.QueryRow(getUserByNickname, user.Nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	log.Println(err)

	if err == pgx.ErrNoRows {
		return ErrorUserNotFound
	}
	return nil
}
