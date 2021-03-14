package receiver

import (
	"fmt"
	"net/http"

	"github.com/barathvk/sender-receiver/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setupServer(logger *log.Entry) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.POST("/count", func(context *gin.Context) {
		var count common.Count
		context.BindJSON(&count)
		logger.WithFields(log.Fields{"nodeId": count.NodeId, "senderAppId": count.AppId}).Info("Received count ", count.Count)
		context.JSON(http.StatusAccepted, gin.H{"status": "accepted", "count": count.Count})
	})
	return server
}

func Start(appId string, port int) {
	logger := log.WithField("appId", appId)
	logger.Info("Starting receiver...")
	server := setupServer(logger)
	logger.Info("Starting receiver server on port ", port)
	server.Run(fmt.Sprintf(":%d", port))
}
