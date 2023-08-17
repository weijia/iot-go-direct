package iot_go_test

import (
	"iot_go/pkg/thingsboard"
	"iot_go/pkg/thingsboard_shared"
	"testing"
)

func TestCreateDevice(t *testing.T) {
	server := thingsboard.ThingsboardServer{
		DeviceProfile: thingsboard_shared.DeviceProfile{
			ThingsboardServerInfo: thingsboard_shared.ThingsboardServerInfo{
				Server: "http://120.55.92.168",
				Port:   8080,
			},
			ProvisionKey:    "0hsh1hpc605g4kwyal46",
			ProvisionSecret: "68rsgqafhw0anhcwnccr",
		},
	}
	server.CreateDevice("hello")
}

func TestTelemetry(t *testing.T) {
	server := thingsboard.ThingsboardServer{
		DeviceProfile: thingsboard_shared.DeviceProfile{
			ThingsboardServerInfo: thingsboard_shared.ThingsboardServerInfo{
				Server: "http://120.55.92.168",
				Port:   8080,
			},
			ProvisionKey:    "0hsh1hpc605g4kwyal46",
			ProvisionSecret: "68rsgqafhw0anhcwnccr",
		},
	}
	data := map[string]interface{}{
		"color": "13",
	}
	server.UploadTelemetry("hello", data)
}

func TestSubscribe(t *testing.T) {
	server := thingsboard.ThingsboardServer{
		DeviceProfile: thingsboard_shared.DeviceProfile{
			ThingsboardServerInfo: thingsboard_shared.ThingsboardServerInfo{
				Server: "http://120.55.92.168",
				Port:   8080,
			},
			ProvisionKey:    "0hsh1hpc605g4kwyal46",
			ProvisionSecret: "68rsgqafhw0anhcwnccr",
		},
	}
	server.SubscribeToAttribute("hello1", 100)
}
