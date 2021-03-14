package sender

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/barathvk/sender-receiver/common"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

func sendRequest(count int, appId string, nodeId string, logger *log.Entry) {
	logger.Debug("Sending counter with value ", count)
	countPayload := &common.Count{Count: count, AppId: appId, NodeId: nodeId}
	countJson, _ := json.Marshal(countPayload)
	_, err := http.Post("http://localhost:8080/count", "application/json", bytes.NewBuffer(countJson))
	if err != nil {
		logger.Error("Request failed with count ", count, ": ", err)
	} else {
		logger.Info("Sent request with count ", count)
	}
}

func Start(appId string, initialValue int) {
	nodeId := ksuid.New().String()
	logger := log.WithFields(log.Fields{"appId": appId, "nodeId": nodeId})
	logger.Info("starting sender...")
	count := initialValue
	sendRequest(count, appId, nodeId, logger)
	count += 1
	for range time.Tick(time.Second * 1) {
		sendRequest(count, appId, nodeId, logger)
		count += 1
	}
}
