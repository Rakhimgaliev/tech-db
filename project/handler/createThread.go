package handler

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateThread(context *gin.Context) {
	thread := &models.Thread{}

	context.BindJSON(thread)
	thread.Slug = context.Param("slug")

	err := db.CreateThread(h.conn, thread)

	if err != nil {
		context.JSON(500, "application/json")
	}

	threadJSON, _ := json.Marshal(thread)
	context.Data(201, "application/json", threadJSON)
}
