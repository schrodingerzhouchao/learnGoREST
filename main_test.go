package main

import (
	"fmt"
	"gopro/go-rest/routers/api"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/gin-gonic/gin"
)

// Setup build the server for testing
func Setup() *gin.Engine {
	host := "127.0.0.1:27017"
	var co = api.NewConn
	api.DBbuild("appdb", "movies")
	co.Use(host, co.DBName, co.CollName)
	r := gin.Default()
	r.GET("/v1/movies/:movie_id", api.Get)
	return r
}

// TestGet testing the api Get
func TestGet(t *testing.T) {
	testr := Setup()
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/v1/movies/5c9198c70876b7fd4536e44f", nil)
	testr.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Request code is %v", writer.Code)
	} else {
		fmt.Println(writer.Body)
	}
}

func BenchmarkGet(b *testing.B) {

	fcpu, err := os.OpenFile("./cpuinfo.out", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("open file:%v\n", err.Error())
	}
	defer fcpu.Close()
	pprof.StartCPUProfile(fcpu)
	defer pprof.StopCPUProfile()

	fmem, err := os.OpenFile("./meminfo.out", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("open file:%v\n", err.Error())
	}
	defer fmem.Close()
	pprof.WriteHeapProfile(fmem)

	testr := Setup()
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/v1/movies/5c91fce8c2c289313fb3deec", nil)
	for i := 0; i < b.N; i++ {
		testr.ServeHTTP(writer, request)
	}

}
