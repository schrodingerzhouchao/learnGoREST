package db

import (
	"fmt"
	"gopro/go-rest/models"
	"log"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

// DBconn mgo session
type DBconn struct {
	Session    *mgo.Session
	Collection *mgo.Collection
	DBName     string
	CollName   string
}

// Use database
func (conn *DBconn) Use(host, dbName, collName string) {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatalln("connect MongoDB failed,", err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(dbName).C(collName)
	conn.Session = session
	conn.Collection = c
}

// GetDBAndCollName get the connected DB name
func (conn *DBconn) GetDBAndCollName() (map[string][]string, error) {
	dbs, err := conn.Session.DatabaseNames()
	result := make(map[string][]string)
	for _, dbname := range dbs {
		collname, _ := conn.Session.DB(dbname).CollectionNames()
		result[dbname] = collname
	}
	return result, err
}

// Clone copy a new session
//func (conn *DBconn) Clone(dbName, collName string) {
//	newSess := conn.session.Copy()
//	conn.session = newSess
//	conn.collection = newSess.DB(dbName).C(collName)
//}

// Close database
func (conn *DBconn) Close() {
	conn.Session.Close()
}

// GetByID get the movie struct
func (conn *DBconn) GetByID(id string) (models.Movie, error) {
	c := conn.Session.Copy()
	defer c.Close()
	cc := c.DB(conn.DBName).C(conn.CollName)
	var movie models.Movie
	err := cc.FindId(bson.ObjectIdHex(id)).One(&movie)

	if err != nil {
		log.Println(err.Error())
		return movie, err
	}
	return movie, nil
}

// PostStruct post new data
func (conn *DBconn) PostStruct(movie models.Movie) error {
	c := conn.Session.Copy()
	defer c.Close()
	cc := c.DB(conn.DBName).C(conn.CollName)
	err := cc.Insert(movie)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

// DeleteByID delete movie by id_
func (conn *DBconn) DeleteByID(id string) error {
	c := conn.Session.Copy()
	defer c.Close()
	cc := c.DB(conn.DBName).C(conn.CollName)
	err := cc.RemoveId(bson.ObjectIdHex(id))

	return err
}

// UpdateStruct modifies the data by id_
func (conn *DBconn) UpdateStruct(id string, movie models.Movie) error {
	c := conn.Session.Copy()
	defer c.Close()
	cc := c.DB(conn.DBName).C(conn.CollName)
	err := cc.Update(bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"$set": bson.M{
			"name":      movie.Name,
			"year":      movie.Year,
			"directors": movie.Directors,
			"writers":   movie.Writers,
		}})

	return err
}

// UpsertStruct update or add
func (conn *DBconn) UpsertStruct(movie models.Movie) error {
	c := conn.Session.Copy()
	defer c.Close()
	cc := c.DB(conn.DBName).C(conn.CollName)
	// TODO selector as paramter
	selector := bson.M{"name": movie.Name, "year": movie.Year}
	data := bson.M{"$set": bson.M{"name": movie.Name,
		"year":      movie.Year,
		"directors": movie.Directors,
		"writers":   movie.Writers,
	}}
	changeInfo, err := cc.Upsert(selector, data)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("%+v\n", changeInfo)
	return err
}
