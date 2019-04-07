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
			log.Println(err)
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
	// log.Println(user)

	// err := db.GetUserByNickname(h.conn, user)
	// if err != nil {
	// 	if err == db.ErrorUserNotFound {
	// 		context.JSON(404, err)
	// 		return
	// 	}
	// 	context.JSON(500, err)
	// 	return
	// }

}
