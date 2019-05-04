package main

// go install ./project/main/ && ../../../../bin/main
// curl -d --header "Content-Type: application/json" --request POST http://localhost:5000/api/forum/create --data '{"title":"Pirate","user":"j.sparrow","slug":"pirate-stories"}'
// curl --header "Content-Type: application/json"   --request POST   --data '{"title":"qsweqwe","user":"asdaads","slug":"fgdfg"}'   http://localhost:5000/api/forum/create
// curl --header "Content-Type: application/json"   --request POST   --data '{"fullname":"qsweqwe","about":"asdaads","email":"jashhasd@mail.ru"}'   http://localhost:5000/api/user/teuumikr/create

import (
	"github.com/Rakhimgaliev/tech-db-forum/project/handler"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

var dbConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "docker",
	User:     "docker",
	Password: "docker",
}

func main() {
	router := gin.Default()

	handler := handler.NewConnPool(&dbConfig)

	router.POST("/api/forum/:slug", func(c *gin.Context) {
		if c.Param("slug") == "create" {
			handler.CreateForum(c)
		}
	})
	router.GET("/api/forum/:slug/details", handler.GetForum)
	router.GET("/api/forum/:slug/users", handler.GetForumUsers)

	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)

	router.POST("/api/forum/:slug/create", handler.CreateThread)
	router.GET("/api/forum/:slug/threads", handler.GetThreads)
	router.GET("/api/thread/:slug_or_id/details", handler.ThreadDetails)
	router.POST("/api/thread/:slug_or_id/details", handler.ThreadUpdate)
	router.POST("/api/thread/:slug_or_id/vote", handler.CreateThreadVote)

	router.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	router.GET("/api/thread/:slug_or_id/posts", handler.GetPosts)
	router.POST("/api/post/:id/details", handler.UpdatePost)
	router.GET("/api/post/:id/details", handler.GetPost)

	router.Run("127.0.0.1:5000")
}
