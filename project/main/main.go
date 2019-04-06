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

	router.POST("/api/user/:nickname/create", handler.CreateUser)
	router.GET("/api/user/:nickname/profile", handler.GetUser)
	router.POST("/api/user/:nickname/profile", handler.UpdateUser)

	router.POST("/api/forum/:slug/create", handler.CreateThread)

	// router.POST("/forum/", func(c *gin.Context) {
	// 	if c.Param("slug") == "create" {
	// 		c.JSON(200, gin.H{
	// 			"description": `Создание нового форума.`,
	// 		})
	// 	}
	// 	c.JSON(200, gin.H{
	// 		"description": `hz`,
	// 	})
	// })

	// router.POST("/forum/:slug/create", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Добавление новой ветки обсуждения на форум.`,
	// 	})
	// })

	// router.GET("/forum/:slug/users", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
	// 							Пользователи выводятся отсортированные по nickname в порядке возрастания.
	// 							Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.`,
	// 	})
	// })

	// router.GET("/forum/:slug/threads", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение списка ветвей обсужления данного форума.
	// 							Ветви обсуждения выводятся отсортированные по дате создания.`,
	// 	})
	// })

	// router.GET("/post/:id/details", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение информации о ветке обсуждения по его имени.`,
	// 	})
	// })

	// router.POST("/post/:id/details", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Изменение сообщения на форуме.
	// 							Если сообщение поменяло текст, то оно должно получить отметку "isEdited".`,
	// 	})
	// })

	// router.POST("/service/clear", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Очистка всех данных в базе`,
	// 	})
	// })

	// router.GET("/service/status", func(c *gin.Context) {
	// 	Status(c, conn)
	// })

	// router.POST("/thread/:slug_or_id/create", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Добавление новых постов в ветку обсуждения на форум.
	// 							Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).`,
	// 	})
	// })

	// router.GET("/thread/:slug_or_id/details", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение информации о ветке обсуждения по его имени.`,
	// 	})
	// })

	// router.POST("/thread/:slug_or_id/details", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Обновление ветки обсуждения на форуме.`,
	// 	})
	// })

	// router.GET("/thread/:slug_or_id/posts", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение информации о ветке обсуждения по его имени.`,
	// 	})
	// })

	// router.POST("/thread/:slug_or_id/vote", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Изменение голоса за ветвь обсуждения.
	// 							Один пользователь учитывается только один раз и может изменить своё
	// 							мнение.`,
	// 	})
	// })

	// router.POST("/user/:nickname/create", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Создание нового пользователя в базе данных.`,
	// 	})
	// })

	// router.GET("/user/:nickname/profile", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Получение информации о пользователе форума по его имени.`,
	// 	})
	// })

	// router.POST("/user/:nickname/profile", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"description": `Изменение информации в профиле пользователя.`,
	// 	})
	// })

	router.Run("127.0.0.1:5000")
}
