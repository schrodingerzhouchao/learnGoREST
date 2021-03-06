package db

import (
	"fmt"
	"gopro/go-rest/models"
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var conn = DBconn{
	Session:    nil,
	Collection: nil,
	DBName:     "appdb",
	CollName:   "movies",
}

// TestMovie for testing
var TestMovie = models.Movie{
	ID:        bson.NewObjectId(),
	Name:      "The Shawshank Redemption",
	Year:      "1994",
	Directors: []string{"Frank Darabont"},
	Writers:   []string{"Frank Darabont", "Stephen Edwin King"},
}

var TestMovieupdate = models.Movie{
	ID:        bson.NewObjectId(),
	Name:      "The Shawshank Redemption",
	Year:      "2004",
	Directors: []string{"Frank Darabont"},
	Writers:   []string{"Frank Darabont", "Stephen Edwin King"},
}

const host = "127.0.0.1:37017"

// TestUse test the connection
func TestUse(t *testing.T) {
	//host := "127.0.0.1:37017"
	//dbName := "appdb"
	//collName := "movies"

	conn.Use(host, conn.DBName, conn.CollName)
}

// TestGetDBAndCollName test get the collection names
func TestGetDBAndCollName(t *testing.T) {
	result, err := conn.GetDBAndCollName()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(result)
}

// TestGet test the Get method
func TestGet(t *testing.T) {
	id := "5ca76aaec14841c7532f4c04"
	movie, err := conn.GetByID(id)
	if err != nil {
		t.Fatal("Get movie failed:", err)
	}
	if movie.Name != "The Dark Knight" {
		t.Fatal("Get wrong name")
	}
	fmt.Println(movie)
}

func TestGetByYear(t *testing.T) {
	year := "2009"
	conn.Use(host, conn.DBName, conn.CollName)
	movies, err := conn.GetByYear(year)
	if err != nil {
		t.Fatal("Get movie failed:", err)
	}

	fmt.Println(movies)
}

// TestPost test the PostStruct method
func TestPost(t *testing.T) {
	err := conn.PostStruct(TestMovie)
	if err != nil {
		t.Fatal("Post movie failed:", err)
	}
}

func TestDeleteByID(t *testing.T) {
	err := conn.DeleteByID(TestMovie.ID.Hex())
	if err != nil {
		log.Fatal("delete error:", err)
	}
}

// TestUpdate test the Update method
/*
func TestUpdate(t *testing.T) {
	id := "5ca76aaec14841c7532f4c04"
	err := conn.UpdateStruct(id, TestMovieupdate)
	if err != nil {
		t.Fatal("Update movie failed:", err)
	}
	movie, err := conn.GetByID(id)
	fmt.Println("updated movie:", movie)
}
*/
// TestUpsert test upsert method
func TestUpsert(t *testing.T) {
	err := conn.UpsertStruct(TestMovieupdate)
	if err != nil {
		t.Fatal("Upsert movie:", err)
	}
}

func BenchmarkGet(b *testing.B) {
	id := "5ca76aaec14841c7532f4c04"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		movie, err := conn.GetByID(id)
		if err != nil {
			log.Println("Get get error", err)
		} else {
			fmt.Println(i, ":", movie)
		}
	}
}
