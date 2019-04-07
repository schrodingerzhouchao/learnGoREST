package logging

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoggerConn used for logger DB
type LoggerConn struct {
	Session    *mgo.Session
	Collection *mgo.Collection
	DBName     string
	CollName   string
}

// LoggerMessage logger
type LoggerMessage struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	LogType    string        `json:"logtype" bson:"logtype"`
	Time       string        `json:"time" bson:"time"`
	FuncFile   string        `json:"funcfile" bson:"funcfile"`
	Func       string        `json:"func" bson:"func"`
	Line       string        `json:"line" bson:"line"`
	FatalError string        `json:"fatalerror" bson:"fatalerror"`
	Text       RouterText    //`json:"routertext" bson:"routertext"`

}

// RouterText router info
type RouterText struct {
	StatusCode string `json:"statuscode" bson:"statuscode"`
	Latency    string `json:"latency" bson:"latency"`
	ReqIP      string `json:"reqip" bson:"reqip"`
	Method     string `json:"method" bson:"method"`
	Path       string `json:"path" bson:"path"`
}

const (
	//TypeInfo info
	TypeInfo = "INFO"
	//TypeError error
	TypeError = "ERROR"
	//TypeFatal fatal
	TypeFatal = "FATAL"
)

// CreateLoggerMessage create logger struct
func (logm *LoggerMessage) CreateLoggerMessage(logtype string, skip int) LoggerMessage {
	pc, fn, line, _ := runtime.Caller(skip)
	logtime := time.Now()
	logtimeStr := fmt.Sprintf("%d/%02d/%02d-%02d:%02d:%02d",
		logtime.Year(), logtime.Month(), logtime.Day(),
		logtime.Hour(), logtime.Minute(), logtime.Second())
	return LoggerMessage{
		LogType:  logtype,
		Time:     logtimeStr,
		FuncFile: fn,
		Func:     runtime.FuncForPC(pc).Name(),
		Line:     strconv.Itoa(line),
	}
}

// Use connect mongodb
func (logconn *LoggerConn) Use(host, dbName, collName string) error {
	session, err := mgo.Dial(host)
	if err != nil {
		log.Println("connect MongoDB failed", err)
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(dbName).C(collName)
	logconn.Session = session
	logconn.Collection = c

	return nil
}

// Insert logger message into DB
func (logconn *LoggerConn) Insert(loggermessage LoggerMessage) error {
	c := logconn.Session.Copy()
	defer c.Close()
	cc := c.DB(logconn.DBName).C(logconn.CollName)
	err := cc.Insert(loggermessage)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

// GetRouterInfo get statuscode etc.
func GetRouterInfo(c *gin.Context) RouterText {
	start := time.Now()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)
	code := strconv.Itoa(c.Writer.Status())
	path := c.Request.URL.Path
	method := c.Request.Method
	reqip := c.ClientIP()
	return RouterText{
		StatusCode: code,
		Latency:    latency.String(),
		ReqIP:      reqip,
		Method:     method,
		Path:       path,
	}

}

var logm = new(LoggerMessage)

// Logger method
func (logconn *LoggerConn) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logLevel string
		/*
			switch c.Writer.Status() {
			case 200:
				logLevel = TypeInfo
			default:
				logLevel = TypeError
			}
		*/
		if c.Writer.Status() <= 206 && c.Writer.Status() >= 200 {
			logLevel = TypeInfo
		} else {
			logLevel = TypeError
		}
		tmp := logm.CreateLoggerMessage(logLevel, 1)
		tmp.Text = GetRouterInfo(c)
		logconn.Insert(tmp)
	}
}

// Logger2 for fatal
func (logconn *LoggerConn) Logger2(err error) {
	tmp := logm.CreateLoggerMessage(TypeFatal, 1)
	tmp.FatalError = err.Error()
	logconn.Insert(tmp)
}

/*
const (
	host     = "127.0.0.1:37017"
	dbname   = "appdb"
	collname = "moviesLog"
)
*/
