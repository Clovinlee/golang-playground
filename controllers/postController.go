package controllers

import (
	"chris/gochris/initializers"
	"chris/gochris/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {

	var body struct {
		Body  string
		Title string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Body are required"})
		return
	}

	if body.Body == "" || body.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Body are required"})
		return
	}

	post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func PostsIndex(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func PostsShow(c *gin.Context) {
	// ID from URL param
	id := c.Param("id")

	var post models.Post
	// cant do post == nil or whatever
	// because default value of type is not nil ( for example, if int then 0 etc)
	// isntead, can do post == models.Post{} to check if all are default value

	result := initializers.DB.Find(&post, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func PostsUpdate(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	if body.Body == "" || body.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Body are required", "body": body})
		return
	}

	var post models.Post
	result := initializers.DB.First(&post, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func PostsDelete(c *gin.Context) {
	id := c.Param("id")

	result := initializers.DB.Delete(&models.Post{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Record deleted successfully",
	})
}
