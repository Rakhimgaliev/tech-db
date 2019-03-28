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

	err := context.BindJSON(&user)
	if err != nil {
		log.Println(err)
		return
	}
	(*user).Nickname = context.Param("nickname")

	err = db.CreateUser(h.conn, user)
	if err != nil {
		if err == db.ErrorUserAlreadyExists {
			err := db.GetUserByNickname(h.conn, user)
			if err != nil {
				context.JSON(500, err)
				return
			}
			userJSON, _ := json.Marshal(user)
			context.JSON(409, string(userJSON))
			return
		} else {
			context.JSON(500, err)
			return
		}
	}

	userJSON, _ := json.Marshal(user)
	context.JSON(200, string(userJSON))
}

func (h handler) GetUserBySlug(context *gin.Context) {
	user := &models.User{}

	(*user).Nickname = context.Param("nickname")

	err := db.GetUserByNickname(h.conn, user)
	if err != nil {
		context.JSON(500, err)
		return
	}

	userJSON, _ := json.Marshal(user)
	context.JSON(200, string(userJSON))
}
