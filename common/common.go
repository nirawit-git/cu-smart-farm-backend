package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

// MqttClient : this is MQTT client that listen to server
var MqttClient *service.Client

// CheckErr : return true and print if error
func CheckErr(source string, err error) bool {
	if err != nil {
		fmt.Printf("\n[ERR] IN %s || %s\n", source, err)
		return true
	}
	return false
}

// ConnectToMQTT : connects to mqtt server and return error if error
func ConnectToMQTT() error {
	if MqttClient != nil {
		MqttClient.Disconnect()
	}
	MqttClient = &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetUsername([]byte("admin"))
	msg.SetPassword([]byte("iyddyoot"))
	msg.SetWillQos(2)
	msg.SetKeepAlive(120)
	msg.SetVersion(3)
	msg.SetClientId([]byte("smart-farm-backend"))
	msg.SetWillTopic([]byte("backend-service"))
	msg.SetWillMessage([]byte("backend: connecting.."))

	// msg.SetCleanSession(true)
	return MqttClient.Connect("tcp://164.115.27.177:1883", msg)
}

// ParseJSON : parse byte to json (gin.H)
func ParseJSON(payload []byte) gin.H {
	var jsonData gin.H
	CheckErr("Parsing JSON", json.Unmarshal(payload, &jsonData))
	return jsonData
}

// ConnectToInfluxDB : connect to influx DB and return client
func ConnectToInfluxDB() (client.Client, error) {
	influxConn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "admin",
		Password: "4fs,mg-0zv",
	})
	CheckErr("Influx Connection", err)
	return influxConn, err
}

// QueryInfluxDB : runs query in influxDB
func QueryInfluxDB(query string) []client.Result {
	clnt, err := ConnectToInfluxDB()

	CheckErr("connect to query influx", err)
	if err == nil {
		resp, err := clnt.Query(client.Query{
			Command:  query,
			Database: "CUSmartFarm",
		})
		CheckErr("querying influx", err)
		if err == nil {
			fmt.Printf("Query Success: %v \n", resp)
		}
		return resp.Results
	}
	return []client.Result{}
}

// WriteInfluxDB : Write a data point in to influxDB
func WriteInfluxDB(measurement string, tags map[string]string, fields map[string]interface{}) error {
	clnt, err := ConnectToInfluxDB()

	CheckErr("connect to query influx", err)
	if err == nil {
		batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "CUSmartFarm",
			Precision: "ms",
		})
		if CheckErr("create batch points", err) {
			return err
		}

		point, err := client.NewPoint("air_sensor", tags, fields, time.Now())
		if CheckErr("create new point", err) {
			return err
		}

		batchPoints.AddPoint(point)
		err = clnt.Write(batchPoints)
		if !CheckErr("querying influx", err) {
			fmt.Printf("DB Write Succeeded\n")
		}
		return nil
	}
	return nil
}
