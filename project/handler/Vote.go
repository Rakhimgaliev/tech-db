package handler

import (
	"encoding/json"

	"github.com/Rakhimgaliev/tech-db/project/db"
	"github.com/Rakhimgaliev/tech-db/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateThreadVote(context *gin.Context) {
	var thread models.Thread
	var vote models.Vote

	err := BindJSON(context, &vote)
	if err != nil {
		context.JSON(400, err)
		return
	}

	err = db.CreateThreadVote(h.conn, context.Param("slug_or_id"), &thread, &vote)
	if err != nil {
		switch err {
		case db.ErrorThreadNotFound:
			context.JSON(404, err)
			return
		default:
			context.JSON(500, err)
		}
		return
	}
	threadJSON, _ := json.Marshal(thread)
	context.Data(200, "application/json", threadJSON)
	return
}
