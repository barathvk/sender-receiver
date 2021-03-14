package receiver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barathvk/sender-receiver/common"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestReceiver(t *testing.T) {
	logger, hook := test.NewNullLogger()
	appId := "test-app"
	nodeId := "test-node"
	router := setupServer(logger.WithField("appId", appId))
	recorder := httptest.NewRecorder()
	countPayload := &common.Count{Count: 4, AppId: appId, NodeId: nodeId}
	countJson, _ := json.Marshal(countPayload)
	req, _ := http.NewRequest("POST", "/count", bytes.NewBuffer(countJson))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusAccepted, recorder.Code)
	var response map[string]interface{}
	json.NewDecoder(recorder.Body).Decode(&response)
	assert.Equal(t, http.StatusAccepted, recorder.Code)
	assert.Equal(t, "accepted", response["status"])
	assert.Equal(t, "Received count 4", hook.LastEntry().Message)
}
