package handler

import (
	"encoding/json"
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
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
	log.Println("HERE: ", err)
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
	context.Data(201, "application/json", threadJSON)
	return
}
