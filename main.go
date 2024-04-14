package main

import (
	"chris/gochris/controllers"
	"chris/gochris/initializers"
	"chris/gochris/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/register", controllers.UserRegister)
	r.POST("/login", controllers.UserLogin)

	postGroup := r.Group("/posts").Use(middleware.RequireAuth)
	{
		postGroup.GET("/", controllers.PostsIndex)
		postGroup.POST("/", controllers.PostsCreate)
		postGroup.GET("/:id", controllers.PostsShow)
		postGroup.PUT("/:id", controllers.PostsUpdate)
		postGroup.DELETE("/:id", controllers.PostsDelete)
	}

	r.GET("/cookie", func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "No cookie found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"cookie": tokenString,
			})
		}

	})

	r.Run() // listen and serve based on port of env
}
