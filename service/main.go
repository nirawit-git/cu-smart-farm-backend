package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/config"

	"github.com/rod41732/cu-smart-farm-backend/common"

	"github.com/rod41732/cu-smart-farm-backend/service/receiver"
	"github.com/rod41732/cu-smart-farm-backend/service/worker"
)

func main() {
	config.Init()
//	print("[Debug] target: " + config.MQTT["address"])
	common.ShouldPrintDebug = true
	trigger := new(receiver.Trigger)
	rpc.Register(trigger)
	rpc.HandleHTTP()

	// common.InitializeKeyPair()
	// middleware.Initialize()

	time.Sleep(2 * time.Second) // wait until MQTT connect
	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	worker.Init()
	go worker.Work()
	http.Serve(l, nil)
}
