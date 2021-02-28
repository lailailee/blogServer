package router

import (
	"end/controllers"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func GinRouter(router *gin.Engine) {
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "OPTIONS, GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type,Company",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))
	router.GET("/api/tags", controllers.GetTag)
	router.POST("/api/tag", controllers.AddTag)
	router.PUT("/api/tag/id", controllers.EditTag)
	router.DELETE("/api/tag/id", controllers.DeleteTag)

	router.GET("/api/categorys", controllers.GetCategory)
	router.POST("/api/category", controllers.AddCategory)
	router.PUT("/api/category/:id", controllers.EditCategory)
	router.DELETE("/api/category/:id", controllers.DeleteCategory)

	router.GET("/api/articles", controllers.GetArticles)
	router.GET("/api/article/:id", controllers.GetArticle)
	router.POST("/api/article", controllers.AddArticle)
	router.PUT("/api/article/id", controllers.EditArticle)
	router.DELETE("/api/article/id", controllers.DeleteArticle)
}
