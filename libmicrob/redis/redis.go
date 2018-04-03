package redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

var pool *redis.Pool
var host string
var redisDb int
var Hostname string

func InitRedis(conf *types.Conf) *terr.Trace {
	msgs.Status("Initializing Redis connection")
	host = conf.RedisAddr
	redisDb = conf.RedisDb
	pool = newPool(conf.RedisAddr)
	Hostname = conf.Name
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		tr := terr.New("services.logs.redis.initRedis", err)
		return tr
	}
	return nil
}

func GetConn() redis.Conn {
	conn := pool.Get()
	conn.Do("SELECT", redisDb)
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

func GetKeys(key string) ([]interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	var data interface{}
	data, err := conn.Do("LRANGE", key, 0, -1)
	var keys []interface{}
	if err != nil {
		return keys, fmt.Errorf("error getting key %s: %v", key, err)
	}
	// delete the list
	_, err = conn.Do("DEL", key)
	if err != nil {
		return keys, errors.New("Can not delete key " + key)
	}
	keys = data.([]interface{})
	return keys, err
}

func GetKey(skey string) (interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	var data interface{}
	data, err := conn.Do("GET", skey)
	var key interface{}
	if err != nil {
		//msg := "Can not get key " + skey + " from Redis"
		tr := terr.New("services.logs.redis.Set", err)
		err = tr.ToErr()
		return key, fmt.Errorf("error getting key %s: %v", skey, err)
	}
	// delete the list
	_, err = conn.Do("DEL", skey)
	if err != nil {
		return key, errors.New("Can not delete key " + skey)
	}
	key = data.(interface{})
	return key, err
}

func PushKey(key string, value []byte) error {
	conn := GetConn()
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
