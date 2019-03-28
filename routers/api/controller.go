package api

import (
	"fmt"
	"gopro/go-rest/db"
	"gopro/go-rest/models"

	"github.com/gin-gonic/gin"
)

var movie = new(models.Movie)
var NewConn = new(db.DBconn)

func DBbuild(dbname, collname string) {
	NewConn.DBName = dbname
	NewConn.CollName = collname
}

// Get gets the movie info
func Get(c *gin.Context) {
	id := c.Param("movie_id")
	fmt.Println(id)
	movie, err := NewConn.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	} else {
		c.JSON(200, gin.H{"info": movie})
	}

}

// Create adds the movie
func Create(c *gin.Context) {

}

// Update modify the movie
func Update(c *gin.Context) {

}

// Upsert update or add
func Upsert(c *gin.Context) {

}

// Abc test
func Abc(c *gin.Context) {
	id := c.Param("abc_id")
	c.JSON(200, gin.H{"message": id})
}
