package handler

import (
	"encoding/json"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateUser(context *gin.Context) {
	user := &models.User{}

	err := context.BindJSON(user)
	if err != nil {
		log.Println(err)
		return
	}
	(*user).Nickname = context.Param("nickname")

	err = db.CreateUser(h.conn, user)
	if err != nil {
		if err == db.ErrorUserAlreadyExists {
			err = db.GetUserByEmail(h.conn, user)
			userJSON, _ := json.Marshal(user)
			context.Data(409, "application/json", userJSON)
			return
		} else {
			context.JSON(500, err)
			log.Println(err)
			return
		}
	}

	userJSON, _ := json.Marshal(user)
	context.Data(201, "application/json", userJSON)
}

func (h handler) GetUser(context *gin.Context) {
	user := &models.User{}

	(*user).Nickname = context.Param("nickname")

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
