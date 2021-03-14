package receiver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/barathvk/sender-receiver/common"
	"github.com/gin-gonic/gin"
	cmd "github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
)

var lastRequest time.Time
var latestCount int
var currentSender *cmd.Cmd

func startSender(appId string, initialCount int) {
	log.Info("Starting sender with initial count ", initialCount)
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	currentSender = cmd.NewCmdOptions(cmdOptions, "./sender-receiver", "--sender", "--initial-count", strconv.Itoa(latestCount+1))
	<-currentSender.Start()
}

func heartbeat(appId string) {
	for range time.Tick(time.Millisecond * 1000) {
		if time.Since(lastRequest) > time.Second*1 {
			log.Warn("Sender is not running")
			startSender(appId, latestCount)
		}
	}
}

func setupServer(logger *log.Entry) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.POST("/count", func(context *gin.Context) {
		var count common.Count
		context.BindJSON(&count)
		lastRequest = time.Now()
		latestCount = count.Count
		logger.WithFields(log.Fields{"nodeId": count.NodeId, "senderAppId": count.AppId}).Info("Received count ", count.Count)
		context.JSON(http.StatusAccepted, gin.H{"status": "accepted", "count": count.Count})
	})
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

func Start(appId string, port int) {
	logger := log.WithField("appId", appId)
	logger.Info("Starting receiver...")
	server := setupServer(logger)
	go heartbeat(appId)
	logger.Info("Starting receiver server on port ", port)
	err := server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(err)
	}
}
