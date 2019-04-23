package handler

import (
	"encoding/json"
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateThread(context *gin.Context) {
	thread := &models.Thread{}

	context.BindJSON(thread)
	thread.Forum = context.Param("slug")

	err := db.CreateThread(h.conn, thread)
	if err != nil {
		switch err {
		case db.ErrorUniqueViolation:
			threadJSON, _ := json.Marshal(thread)
			context.Data(409, "application/json", threadJSON)
			return
		case db.ErrorUserNotFound:
			context.JSON(404, err)
			return
		}
		context.JSON(500, "application/json")
		return
	}

	threadJSON, _ := json.Marshal(thread)
	context.Data(201, "application/json", threadJSON)
}

func (h handler) GetThreads(context *gin.Context) {
	queryArgs := context.Request.URL.Query()
	// log.Print(queryArgs["limit"])
	// log.Print(queryArgs["desc"])
	// log.Print(queryArgs["created"])

	limit := 0
	if len(queryArgs["limit"]) > 0 {
		limit, _ = strconv.Atoi(queryArgs["limit"][0])
	}

	since := ""
	if len(queryArgs["since"]) > 0 {
		since = queryArgs["since"][0]
	}

	desc := false
	if len(queryArgs["desc"]) > 0 {
		if queryArgs["desc"][0] == "true" {
			desc = true
		}
	}

	threads := &models.Threads{}

	err := db.GetThreads(h.conn, context.Param("slug"),
		limit,
		since,
		desc,
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
