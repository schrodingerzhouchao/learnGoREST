package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopro/go-rest/models"
	"gopro/go-rest/routers/api"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// Setup build the server for testing
func Setup(path string, method string) *gin.Engine {
	host := "127.0.0.1:37017"
	var co = api.NewConn
	api.DBbuild("appdb", "movies")
	co.Use(host, co.DBName, co.CollName)
	r := gin.Default()
	switch method {
	case "GET":
		r.GET(path, api.Get)
	case "POST":
		r.POST(path, api.Create)
	default:
		log.Fatalln("Error method")
	}

	return r
}

// TestGet testing the api Get
func TestGet(t *testing.T) {
	testr := Setup("/v1/movies/g/:movie_id", "GET")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/v1/movies/g/5ca76aaec14841c7532f4c04", nil)
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

	testr := Setup("/v1/movies/g/:movie_id", "GET")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/v1/movies/g/5ca76aaec14841c7532f4c04", nil)
	for i := 0; i < b.N; i++ {
		testr.ServeHTTP(writer, request)
	}

}

func TestPost(t *testing.T) {

	//var movie = models.Movie{
	//	ID:        bson.NewObjectId(),
	//	Name:      "testMovie",
	//	Year:      "2019",
	//	Directors: []string{"abc bcd"},
	//	Writers:   []string{"xyz", "uvw"},
	//}
	movie := getRandMovie()
	testJSON, _ := json.Marshal(movie)
	testMovie := bytes.NewBuffer([]byte(testJSON))

	testr := Setup("/v1/movies/c", "POST")
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/v1/movies/c", testMovie)
	request.Header.Add("Content-type", "application/json")

	testr.ServeHTTP(writer, request)
}

func BenchmarkPost1(b *testing.B) {
	host := "127.0.0.1:37017"
	var co = api.NewConn
	api.DBbuild("appdb", "movies")
	co.Use(host, co.DBName, co.CollName)

	r := gin.Default()
	r.POST("/v1/movies/c", api.Create)

	for i := 0; i < b.N; i++ {
		movie := getRandMovie()
		testJSON, _ := json.Marshal(movie)
		testMovie := bytes.NewBuffer([]byte(testJSON))
		//testr := Setup("/v1/movies/c", "POST")
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/v1/movies/c", testMovie)
		request.Header.Add("Content-type", "application/json")

		r.ServeHTTP(writer, request)
	}

}

func getRandMovie() models.Movie {
	movie := models.Movie{
		ID:        bson.NewObjectId(),
		Name:      getRandomString(rand.Intn(12), "STR"),
		Year:      getRandomYear(), //getRandomString(4, "NUM"),
		Directors: []string{getRandomString(rand.Intn(10), "STR")},
		Writers:   []string{getRandomString(rand.Intn(10), "STR"), getRandomString(rand.Intn(10), "STR")},
	}
	return movie
}

func getRandomString(l int, opt string) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//num := "0123456789"
	var bytes []byte
	switch opt {
	case "STR":
		bytes = []byte(str)
	//case "NUM":
	//	bytes = []byte(num)
	default:
		log.Fatalln("error option")
	}

	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func getRandomYear() string {
	rand.Seed(time.Now().Unix())
	year := rand.Int63n(2019-1888) + 1888
	return strconv.Itoa(int(year))
}
