package handler

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db/project/db"
	"github.com/Rakhimgaliev/tech-db/project/models"

	"github.com/gin-gonic/gin"
)

func (h handler) CreateForum(context *gin.Context) {
	forum := &models.Forum{}

	err := context.BindJSON(&forum)
	if err != nil {
		return
	}

	err = db.CreateForum(h.conn, forum)
	if err != nil {
		switch err {
		case db.ErrorUserNotFound:
			context.JSON(404, err)
			return
		case db.ErrorForumAlreadyExists:
			err := db.GetForumBySlug(h.conn, forum)
			if err != nil {
				context.JSON(500, err)
				return
			}
			forumJSON, _ := json.Marshal(forum)
			context.Data(409, "application/json", forumJSON)
			return
		default:
			context.JSON(500, err)
			return
		}
	}

	forumJSON, _ := json.Marshal(forum)
	context.Data(201, "application/json", forumJSON)
}

func (h handler) GetForum(context *gin.Context) {
	forum := &models.Forum{}
	forum.Slug = context.Param("slug")

	err := db.GetForumBySlug(h.conn, forum)
	if err != nil {
		if err == db.ErrorForumNotFound {
			context.JSON(404, err)
			return
		}
		context.JSON(500, err)
		return
	}
	forumJSON, _ := json.Marshal(forum)
	context.Data(200, "application/json", forumJSON)
}
