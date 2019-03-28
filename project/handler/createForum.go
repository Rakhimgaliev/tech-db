package handler

import (
	"encoding/json"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"

	"github.com/gin-gonic/gin"
)

func (h handler) CreateForum(context *gin.Context) {
	forum := &models.Forum{}

	err := context.BindJSON(&forum)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.CreateForum(h.conn, forum)
	if err != nil {
		switch err {
		case db.ErrorUserNotFound:
			context.JSON(404, err)
			return
		case db.ErrorForumAlreadyExist:
			err := db.GetForumBySlug(h.conn, forum)
			if err != nil {
				context.JSON(500, err)
				return
			}
			forumJSON, _ := json.Marshal(forum)
			context.JSON(409, string(forumJSON))
			return
		default:
			context.JSON(500, err)
			return
		}
	}

	forumJSON, _ := json.Marshal(forum)
	context.JSON(200, string(forumJSON))
}
