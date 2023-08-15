package thingsboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iot_go/pkg/util"
	"net/http"

	"iot_go/pkg/thingsboard_shared"
)

type ThingsboardServer struct {
	thingsboard_shared.DeviceProfile
}

func (thingsboardServer ThingsboardServer) Post(url string, data interface{}) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		util.IotLogFatal(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(dataStr))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}

func (thingsboardServer ThingsboardServer) CreateDevice(deviceName string) {
	url := fmt.Sprintf("%s:%d/api/v1/provision",
		thingsboardServer.ThingsboardServerInfo.Server,
		thingsboardServer.ThingsboardServerInfo.Port)
	data := map[string]interface{}{
		"deviceName":            deviceName,
		"provisionDeviceKey":    thingsboardServer.ProvisionKey,
		"provisionDeviceSecret": thingsboardServer.ProvisionSecret,
		"credentialsType":       "ACCESS_TOKEN",
		"token":                 deviceName,
	}
	thingsboardServer.Post(url, data)
}

func (thingsboardServer ThingsboardServer) UploadTelemetry(accessToken string, data interface{}) {
	url := fmt.Sprintf("%s:%d/api/v1/%s/telemetry",
		thingsboardServer.ThingsboardServerInfo.Server,
		thingsboardServer.ThingsboardServerInfo.Port,
		accessToken)
	thingsboardServer.Post(url, data)
}
