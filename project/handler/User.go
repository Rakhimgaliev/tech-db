package handler

import (
	"encoding/json"
	"strconv"

	"github.com/Rakhimgaliev/tech-db/project/db"
	"github.com/Rakhimgaliev/tech-db/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateUser(context *gin.Context) {
	user := &models.User{}

	err := context.BindJSON(user)
	if err != nil {
		return
	}
	user.Nickname = context.Param("nickname")

	err = db.CreateUser(h.conn, user)
	if err != nil {
		if err == db.ErrorUserAlreadyExists {
			users, err := db.GetUserByEmailOrNickname(h.conn, user.Email, user.Nickname)
			if err != nil {
				context.JSON(500, "")
			}
			usersJSON, err := json.Marshal(users)
			if err != nil {
				context.JSON(500, "")
			}

			context.Data(409, "application/json", usersJSON)
			return
		} else {
			context.JSON(500, err)
			return
		}
	}

	userJSON, _ := json.Marshal(user)
	context.Data(201, "application/json", userJSON)
}

func (h handler) GetUser(context *gin.Context) {
	user := &models.User{}
	user.Nickname = context.Param("nickname")

	err := db.GetUserByNickname(h.conn, user)
	if err != nil {
		if err == db.ErrorUserNotFound {
			context.JSON(404, err)
			return
		}
		context.JSON(500, err)
		return
	}

	userJSON, _ := json.Marshal(user)
	context.Data(200, "application/json", userJSON)
}

func (h handler) UpdateUser(context *gin.Context) {
	user := &models.User{}
	user.Nickname = context.Param("nickname")

	updateUser := &models.UserUpdate{}
	err := context.BindJSON(updateUser)
	if err != nil {
		context.JSON(500, err)
		return
	}

	err = db.UpdateUser(h.conn, user, updateUser)
	if err != nil {
		if err == db.ErrorUserNotFound {
			context.JSON(404, err)
			return
		}
		if err == db.ErrorUserAlreadyExists {
			context.JSON(409, err)
			return
		}
	}

	userJSON, _ := json.Marshal(user)
	context.Data(200, "application/json", userJSON)
}

func (h handler) GetForumUsers(context *gin.Context) {
	queryArgs := context.Request.URL.Query()

	limit := 0
	if len(queryArgs["limit"]) > 0 {
		limit, _ = strconv.Atoi(queryArgs["limit"][0])
	}

	since := ""
	if len(queryArgs["since"]) > 0 {
		since = queryArgs["since"][0]
	}

	desc := false
	if len(queryArgs["desc"]) > 0 {
		if queryArgs["desc"][0] == "true" {
			desc = true
		}
	}

	users := models.Users{}
	err := db.GetUsersByForum(h.conn, context.Param("slug"),
		limit, since, desc, &users)

	if err != nil {
		if err == db.ErrorForumNotFound {
			context.JSON(404, err)
			return
		}
		context.JSON(500, err)
		return
	}
	usersJSON, _ := json.Marshal(users)
	context.Data(200, "application/json", usersJSON)
}
