package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"log"
	"strconv"
	"strings"
	"time"
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/synw/microb/conf"
	"github.com/synw/microb/utils"
	"github.com/synw/microb/db/rethinkdb"
)
 

type Hit struct {
	Url *url.URL
	Method string
	Ip string
	User_agent string
	Referer string
	Date time.Time
}

var Config = conf.GetConf()
var HitsKeyName = Config["domain"].(string)+"_hits"

func newPool(addr string) *redis.Pool {
  return &redis.Pool{
    MaxIdle: 3,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
  }
}

var (
  pool *redis.Pool
  redisServer = flag.String("redisServer", ":6379", "")
)

func connect() redis.Conn {
	pool = newPool(*redisServer)
	conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        log.Fatal(err)
    }
    return conn
}

func WriteHit(request *http.Request) {
	Conn := connect()
	defer Conn.Close()
	user_agent := strings.Join(request.Header["User-Agent"], ",")
	referer := "nil"
	val, ok := request.Header["Referer"]
    if (ok == true) {
    	referer = strings.Join(val, ",")
    }
	s := "#!#"
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	hit_str := request.URL.String()+s+request.Method+s+request.RemoteAddr+s+user_agent+s+referer+s+ts
	_, err := Conn.Do("LPUSH", &HitsKeyName, hit_str)
	if err != nil {
        fmt.Println("KEYS: error writing key in Redis:", err)
    }
}

func storeHits(quiet bool) {
	Conn := connect()
	defer Conn.Close()
	// get hits set
	listlen, err := redis.Int(Conn.Do("LLEN", &HitsKeyName))
    if err != nil {
        fmt.Println("KEYS: error retrieving Redis list len:", err)
    }
	now := time.Now()
	date := strconv.Itoa(now.Hour())+":"+strconv.Itoa(now.Minute())+":"+strconv.Itoa(now.Second())
	if listlen > 0 {
		values, err := redis.Strings(Conn.Do("LRANGE", &HitsKeyName, 0, listlen))
		if err != nil {
	        fmt.Println("KEYS: error retrieving Redis list values:", err)
	    }
		// save the keys into the db
		go rethinkdb.SaveHits(values)
		_, err = Conn.Do("DEL", &HitsKeyName)
	    if err != nil {
	        fmt.Println("DEL: error deleting Redis keys:", err)
	    }
	    // then report
	    if (quiet == false) {
	    	fmt.Println(date, "-", listlen, "hits")
	    }
	} else {
		if (quiet == false) {
    		fmt.Println(date, "- 0 hits")
    	}
    }
}

func WatchHits(frequency int)  {
	utils.PrintEvent("info", "watching hits")
	for {
		duration := time.Duration(frequency)*time.Second
		for range time.Tick(duration) {
			go storeHits(true)
		}
	}
}
