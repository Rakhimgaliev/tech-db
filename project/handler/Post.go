package handler

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreatePosts(context *gin.Context) {
	posts := models.Posts{}

	context.BindJSON(posts)

	err := db.CreatePosts(h.db, context.Param("slug_or_id"), &posts)

	postsJSON, _ := json.Marshal(posts)
	context.Data(201, "application/json", postsJSON)
}
