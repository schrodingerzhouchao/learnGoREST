package main

import (
	"gopro/go-rest/logger/logging"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var loggerconn = logging.LoggerConn{
	DBName:   "appdb",
	CollName: "moviesLog",
}

// FatalHandler get some error
func FatalHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := loggerconn.Use("123", "name", "name")
		if err != nil {
			log.Println("NOTE here FatalHandler")
			loggerconn.Logger2(err)
		}
	}

}

// InfoHandler get some info
func InfoHandler() error {
	return nil
}

const host = "127.0.0.1:37017"

func main() {
	logconnerr := loggerconn.Use(host, loggerconn.DBName, loggerconn.CollName)
	log.Println("success?", logconnerr)
	//logm := new(LoggerMessage)

	router := gin.Default()
	router.Use(loggerconn.Logger())
	//router := InitRouter()
	router.GET("/test/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "loggingggg"})
	})
	fh := FatalHandler()
	router.GET("/test/fatal", fh)

	router.Run(":9998")
}
