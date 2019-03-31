package handler

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func sendJsonResponse(context *gin.Context, resp json.Marshaler, status int) {
	jsonResponse, err := resp.MarshalJSON()
	if err != nil {
		context.JSON(500, "")
		log.Println(err)
	}

	context.Data(status, "application/json", jsonResponse)
}
