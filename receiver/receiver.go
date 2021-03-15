package receiver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barathvk/sender-receiver/common"
	"github.com/gin-gonic/gin"
	cmd "github.com/go-cmd/cmd"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var lastRequest time.Time
var currentSender *cmd.Cmd

func startSender(appId string) {
	log.Info("Starting sender...")
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	currentSender = cmd.NewCmdOptions(cmdOptions, "./sender-receiver", "--sender")
	<-currentSender.Start()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			status := currentSender.Status()
			n := len(status.Stdout)
			if n > 0 {
				log.Info(status.Stdout[n-1])
			}

		}
	}()
}

func heartbeat(appId string) {
	for range time.Tick(time.Millisecond * 1000) {
		if time.Since(lastRequest) > time.Second*1 {
			log.Warn("Sender is not running")
			startSender(appId)
		}
	}
}
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
			lastRequest = time.Now()
			logLatestCount(redisConnString, logger)
		case error:
			logger.Error(msg)
		}
	}
}
func setupServer(logger *log.Entry) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.POST("/stop", func(context *gin.Context) {
		status := currentSender.Status()
		logger.Warn("stopping sender with pid ", status.PID, " (runtime ", status.Runtime, ")")
		err := currentSender.Stop()
		if err != nil {
			logger.Warn(err)
			context.JSON(http.StatusInternalServerError, gin.H{"status": "error", "pid": status.PID, "error": err.Error()})
		} else {
			context.JSON(http.StatusAccepted, gin.H{"status": "accepted", "pid": status.PID})
		}
	})
	return server
}

func Start(appId string, port int, redisConnString string) {
	logger := log.WithField("appId", appId)
	server := setupServer(logger)
	go heartbeat(appId)
	go subscribe(redisConnString, logger)
	logger.Info("Redis cluster is at ", redisConnString)
	logger.Info("Starting receiver server on port ", port)
	err := server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(err)
	}
}
