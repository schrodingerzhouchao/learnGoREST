package api

import (
	"gopro/go-rest/db"
	"gopro/go-rest/models"
	"log"

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
	//fmt.Println(id)
	movie, err := NewConn.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
	} else {
		c.JSON(200, gin.H{"info": movie})
	}

}

// Create adds the movie
func Create(c *gin.Context) {
	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.String(406, err.Error())
		c.Abort()
		return
	}

	err := NewConn.PostStruct(movie)
	if err != nil {
		log.Println(err)
	}
}

// Update modify the movie
func Update(c *gin.Context) {
	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.String(406, err.Error())
		c.Abort()
		return
	}
	id := c.Param("movie_id")
	err := NewConn.UpdateStruct(id, movie)
	if err != nil {
		log.Println(err)
	}
}

// Upsert update or add
func Upsert(c *gin.Context) {
	var movie models.Movie
	if err := c.BindJSON(&movie); err != nil {
		c.String(406, err.Error())
		c.Abort()
		return
	}

	err := NewConn.UpsertStruct(movie)
	if err != nil {
		log.Println(err)
	}
}

// Abc test
func Abc(c *gin.Context) {
	id := c.Param("abc_id")
	c.JSON(200, gin.H{"message": id})
}
