package main

import (
	"gin-blog-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDb()

	r := gin.Default()

	r.Run(":8080")
}
