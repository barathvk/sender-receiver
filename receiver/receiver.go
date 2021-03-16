package receiver

import (
	"encoding/json"

	"github.com/barathvk/sender-receiver/common"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func logLatestCount(redisConnString string, logger *log.Entry) {
	conn := common.ConnectToRedis(redisConnString)
	cachedCount, err := redis.String(conn.Do("GET", "count"))
	if err != nil {
		logger.Error(err)
	}
	var count common.Count
	json.Unmarshal([]byte(cachedCount), &count)
	logger.WithFields(log.Fields{"nodeId": count.NodeId, "senderAppId": count.AppId}).Info("Received count ", count.Count)
}

func subscribe(redisConnString string, logger *log.Entry) {
	conn := common.ConnectToRedis(redisConnString)
	conn.Do("CONFIG", "SET", "notify-keyspace-events", "KEA")
	psc := redis.PubSubConn{Conn: conn}
	psc.PSubscribe("__key*__:count")
	for {
		switch msg := psc.Receive().(type) {
		case redis.Message:
			logLatestCount(redisConnString, logger)
		case error:
			logger.Error(msg)
		}
	}
}

func Start(appId string, redisConnString string) {
	logger := log.WithField("appId", appId)
	logger.Info("Redis cluster is at ", redisConnString)
	subscribe(redisConnString, logger)
}
