package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateThread(context *gin.Context) {
	thread := &models.Thread{}

	context.BindJSON(thread)
	log.Print(thread.Slug)

	err := db.CreateThread(h.conn, thread)

	if err != nil {
		context.JSON(500, "application/json")
		return
	}

	threadJSON, _ := json.Marshal(thread)
	context.Data(201, "application/json", threadJSON)
}

func (h handler) GetThreads(context *gin.Context) {
	queryArgs := context.Request.URL.Query()
	log.Print(queryArgs["limit"])
	log.Print(queryArgs["desc"])
	log.Print(queryArgs["created"])

	var limit int
	if len(queryArgs["limit"]) > 0 {
		limit, _ = strconv.Atoi(queryArgs["limit"][0])
	}

	var since string
	if len(queryArgs["created"]) > 0 {
		since = queryArgs["created"][0]
	}

	threads := &models.Threads{}

	err := db.GetThreads(h.conn, context.Param("slug"),
		limit,
		since,
		true,
		threads)

	if err != nil {
		if err == db.ErrorForumNotFound {
			context.JSON(404, err)
			return
		}
		context.JSON(500, err)
		return
	}
	threadsJSON, _ := json.Marshal(threads)
	context.Data(200, "application/json", threadsJSON)
	return
}
