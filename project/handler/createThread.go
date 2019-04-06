package handler

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateThread(context *gin.Context) {
	thread := &models.Thread{}

	err := db.CreateThread(h.conn, thread)

	threadJSON, _ := json.Marshal(thread)
	context.Data(201, "application/json", threadJSON)
}
