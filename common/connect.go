package common

import "github.com/gomodule/redigo/redis"

func ConnectToRedis(redisConnString string) redis.Conn {
	conn, err := redis.Dial("tcp", redisConnString)
	if err != nil {
		panic(err)
	}
	return conn
}
