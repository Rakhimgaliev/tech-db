package handler

import (
	"encoding/json"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindJSON(context *gin.Context, obj interface{}) error {
	if err := binding.JSON.Bind(context.Request, obj); err != nil {
		context.Error(err).SetType(gin.ErrorTypeBind)
		return err
	}
	return nil
}

func (h handler) CreatePosts(context *gin.Context) {
	posts := models.Posts{}

	err := BindJSON(context, &posts)

	err = db.CreatePosts(h.conn, context.Param("slug_or_id"), &posts)
	log.Println(err, posts)
	if err != nil {
		switch err {
		case db.ErrorPostCreateBadRequest:
			postsJSON, _ := json.Marshal(posts)
			context.Data(201, "application/json", postsJSON)
		case db.ErrorThreadNotFound, db.ErrorForumNotFound:
			context.JSON(404, err)
		case db.ErrorPostCreateConflict:
			context.JSON(409, err)
		default:
			context.JSON(500, err)
		}
		return
	}

	postsJSON, _ := json.Marshal(posts)
	context.Data(201, "application/json", postsJSON)
}
