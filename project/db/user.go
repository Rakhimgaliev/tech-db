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

	updateUserCommand = `
		UPDATE "user"
			SET fullname = $2, about = $3, email = $4
			WHERE nickname = $1
			RETURNING nickname, fullname, about, email`
)

func CreateUser(conn *pgx.ConnPool, user *models.User) error {
	err := conn.QueryRow(createUser, user.Nickname, user.Fullname, user.About, user.Email).
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

func UpdateUser(conn *pgx.ConnPool, user *models.User, updateUser *models.UserUpdate) error {
	err := GetUserByNickname(conn, user)
	if err != nil {
		return err
	}

	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	if updateUser.Fullname != "" {
		user.Fullname = updateUser.Fullname
	}

	if updateUser.About != "" {
		user.About = updateUser.About
	}

	err = conn.QueryRow(updateUserCommand, user.Nickname, user.Fullname, user.About, user.Email).
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
