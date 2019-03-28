package handler

import (
	"log"

	"github.com/Rakhimgaliev/tech-db-forum/project/models"

	"github.com/gin-gonic/gin"
)

func (handler) CreateForum(context *gin.Context) {
	forum := models.Forum{}
	err := context.BindJSON(&forum)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(context.ContentType())
	log.Println("Title: ", forum.Title)
	log.Println("User: ", forum.User)
	log.Println("Slug: ", forum.Slug)

}
