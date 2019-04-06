package main

import (
	"gopro/go-rest/db"
	"gopro/go-rest/routers/api"

	"github.com/gin-gonic/gin"
)

func main() {
	var co = api.NewConn
	api.DBbuild("appdb", "movies")
	co.Use(db.Host, co.DBName, co.CollName)

	r := gin.Default()
	r.GET("/v1/movies/g/:movie_id", api.Get)
	r.POST("/v1/movies/c", api.Create)
	r.PUT("/v1/movies/u/:movie_id", api.Update)
	r.PUT("/v1/movies/s", api.Upsert)
	r.Run(":8888")
}
