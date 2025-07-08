package main

import (
	"gin-blog-api/database"
	"gin-blog-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
)

func main() {

	database.ConnectDb()
	database.ConnectRedis()

	r := gin.Default()

	routes.User(r)
	// routes.WebSocketRoutes(router)

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.Run(":8080")
}
