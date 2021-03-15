package sender

import (
	"encoding/json"
	"time"

	"github.com/barathvk/sender-receiver/common"
	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

func sendRequest(count int, appId string, nodeId string, conn redis.Conn, logger *log.Entry) {
	logger.Debug("Sending counter with value ", count)
	countPayload := &common.Count{Count: count, AppId: appId, NodeId: nodeId}
	countJson, _ := json.Marshal(countPayload)
	_, err := conn.Do("SET", "count", string(countJson))
	if err != nil {
		logger.Error("Request failed with count ", count, ": ", err)
	} else {
		logger.Info("Sent request with count ", count)
	}
}

func getInitialCount(conn redis.Conn) int {
	cachedCount, err := redis.String(conn.Do("GET", "count"))
	if err != nil {
		log.Warn(err)
		return 0
	}
	var count common.Count
	json.Unmarshal([]byte(cachedCount), &count)
	return count.Count + 1
}

func Start(appId string, redisConnString string) {
	nodeId := ksuid.New().String()
	logger := log.WithFields(log.Fields{"appId": appId, "nodeId": nodeId})
	conn := common.ConnectToRedis(redisConnString)
	initialCount := getInitialCount(conn)
	logger.Info("starting sender with initial count ", initialCount)
	count := initialCount
	sendRequest(count, appId, nodeId, conn, logger)
	count += 1
	for range time.Tick(time.Second * 1) {
		sendRequest(count, appId, nodeId, conn, logger)
		count += 1
	}
}
