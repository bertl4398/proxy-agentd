package main

import (
  "fmt"
  "log"
  "time"
  "strings"

  "github.com/gomodule/redigo/redis"
)

var RedisConn redis.Conn

func connectRedis(host string, port string) redis.Conn {
  conn, err := redis.Dial("tcp", host+":"+port)
	if err != nil {
	    log.Fatal(err)
	}
  return conn
}

func InitLocalRedis() {
  RedisConn = connectRedis("", "6379")
}

func WriteRedisFlowBatch(flow []string) {
	var do_flush bool
	for _, element := range flow {
		fields := strings.Fields(element)
		if len(fields) == 6 {
			bytes := fields[0]
			frames := fields[1]
			proto := fields[2]
			src := fields[3]
			dst := fields[5]
			t := time.Now().UTC()

			RedisConn.Send("HMSET", fmt.Sprintf("flow/%s/%s", proto, src),
			       				 "bytes", bytes, "frames", frames, "dst", dst, "time", t)
			do_flush = true
		} else {
      log.Printf("Flow not recognized: %s", element)
    }
	}
	if do_flush {
		RedisConn.Flush()
		v, err := RedisConn.Receive()
		if err != nil {
				log.Fatal(err)
		}
		log.Print(v)
	}
}
