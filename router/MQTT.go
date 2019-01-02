package router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	"github.com/rod41732/cu-smart-farm-backend/mqtt"
	"github.com/rod41732/cu-smart-farm-backend/storage"
	"github.com/surgemq/message"
)

func idFromTopic(topic []byte) string {
	return strings.TrimSuffix(strings.TrimPrefix(string(topic), "CUSmartFarm/"), "/svr_recv")
}

// InitMQTT sets handler of mqtt router
func InitMQTT() {
	mqtt.SetHandler(handleMessage)
}

func handleMessage(msg *message.PublishMessage) error {
	topic := string(msg.Topic())
	if strings.HasSuffix(topic, "svr_out") { // skip out message
		return nil
	}

	inMessage := []byte(string(msg.Payload()))
	deviceID := idFromTopic(msg.Topic())
	common.Println("[MQTT] <<< ", string(inMessage))

	var message model.DeviceMessage
	err := json.Unmarshal(inMessage, &message)
	common.Printf("[MQTT] <<< parsed Data=%#v\n", message)

	common.PrintError(err)
	if err == nil {
		device, err := storage.GetDevice(deviceID)
		if common.PrintError(err) && err.Error() != "not found" {
			fmt.Println("  At handleMessage : handleMessage -> GetDevice")
			return err
		}
		common.Printf("[MQTT] --- deviceID=[%s]\n", deviceID)
		user := storage.GetUserStateInfo(device.Owner)
		common.Printf("[MQTT] --- owner=%s\n", device.Owner)
		switch message.Type {
		case "greeting":
			device.BroadCast()
		case "data":
			user.ReportStatus(message.Payload, device.ID)
		}
	} else {
		common.Println("[MQTT] !!! Not a data message")
	}
	return nil
}
