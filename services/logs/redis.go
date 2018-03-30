package logs

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/synw/terr"
	"time"
)

var pool *redis.Pool
var host string
var db int

func initRedis(redisAddr string, redisDb int) *terr.Trace {
	host = redisAddr
	db = redisDb
	pool = newPool(redisAddr)
	conn := getConn()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		tr := terr.New("services.logs.redis.initRedis", err)
		return tr
	}
	return nil
}

func getConn() redis.Conn {
	conn := pool.Get()
	conn.Do("SELECT", db)
	return conn
}

func newPool(redisAddr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func getKeys(key string) ([]byte, error) {
	conn := getConn()
	defer conn.Close()
	var data []byte
	data, err := redis.Bytes(conn.Do("LRANGE", 0, -1))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func setKey(key string, value []byte) error {
	conn := getConn()
	defer conn.Close()
	_, err := conn.Do("RPUSH", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		msg := "Can not set key/value in Redis"
		err := errors.New(msg)
		tr := terr.New("services.logs.redis.Set", err)
		err = tr.ToErr()
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return nil
}
