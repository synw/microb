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
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/centrifugal/gocent"
	"github.com/synw/microb/conf"
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

type NumHits struct {
	Hits int
	Timestamp time.Time
	Domain string
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

func sendHits(num_hits int) {
	secret := Config["centrifugo_secret_key"].(string)
	host := Config["centrifugo_host"].(string)
	port := Config["centrifugo_port"].(string)
	purl := fmt.Sprintf("%s:%s", host, port)
	// connect to Centrifugo
	client := gocent.NewClient(purl, secret, 5*time.Second)
	now := time.Now()
	hit := &NumHits{num_hits, now, Config["domain"].(string)}
	data, err := json.Marshal(hit)
	if err != nil {
	 	println(err.Error())
	 }
	_, err = client.Publish(Config["hits_channel"].(string), data)
	 if err != nil {
	 	println("WS ERROR:", err.Error())
	 }
}

func ProcessHit(request *http.Request,loghit bool, verbosity int, c_display chan string) {
	purl := request.URL.Path
	user_agent := strings.Join(request.Header["User-Agent"], ",")
	var hit_str string
	if (loghit == true) {
		Conn := connect()
		defer Conn.Close()
		referer := "nil"
		val, ok := request.Header["Referer"]
	    if (ok == true) {
	    	referer = strings.Join(val, ",")
	    }
		s := "#!#"
		ts := strconv.FormatInt(time.Now().UnixNano(), 10)
		hit_str = purl+s+request.Method+s+request.RemoteAddr+s+user_agent+s+referer+s+ts
		_, err := Conn.Do("LPUSH", &HitsKeyName, hit_str)
		if err != nil {
	        fmt.Println("KEYS: error writing key in Redis:", err)
	    }
	}
    if (verbosity > 0) {
    	msg := request.Method+" "+purl+" from "+request.RemoteAddr+" - "+user_agent
    	c_display <- msg
    }
}

func storeHits(quiet bool, store_hits bool, monitorhits bool, c chan int) {
	Conn := connect()
	defer Conn.Close()
	// get hits set
	listlen := 0
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
	    if (store_hits == true) {
			// save the keys into the db
			go rethinkdb.SaveHits(values)
		}
		_, err = Conn.Do("DEL", &HitsKeyName)
	    if err != nil {
	        fmt.Println("DEL: error deleting Redis keys:", err)
	    }
	    // then report
	    if (quiet == false) {
	    	fmt.Println(date, "-", listlen, "hits")
	    }
	    if (monitorhits == true) {
	    	sendHits(listlen)
	    }
	} else {
		if (quiet == false) {
    		fmt.Println(date, "- 0 hits")
    	}
    }
    c <- listlen
    return
}

func WatchHits(frequency int, store_hits bool, monitorhits bool, c chan int)  {
	for {
		duration := time.Duration(frequency)*time.Second
		for range time.Tick(duration) {
			go storeHits(true, store_hits, monitorhits, c)
		}
	}
	
}
