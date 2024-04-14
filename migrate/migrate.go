package main

import (
	"chris/gochris/initializers"
	"chris/gochris/models"
	"fmt"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
	initializers.DB.AutoMigrate(&models.User{})
	fmt.Println("Migration Completed")
}
