package storage

import (
	"fmt"

	"github.com/rod41732/cu-smart-farm-backend/model/device"

	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"

	"github.com/rod41732/cu-smart-farm-backend/model/user"
)

// for marshalled db
type userData struct {
	Username string   `json:"username"`
	Devices  []string `json:"devices"`
}

var mappedUserObject = make(map[string]*user.RealUser)

// SetUserStateInfo : Map username into *user.RealUser
func SetUserStateInfo(username string, user *user.RealUser) {
	fmt.Printf("added user: %s\n", username)
	mappedUserObject[username] = user
}

func ensureUser(username string) {
	var tmp userData
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		return
	}
	mdb.DB("CUSmartFarm").C("users").Find(bson.M{
		"username": username,
	}).One(&tmp)
	newUser := user.RealUser{Username: username}
	newUser.Init(tmp.Devices)
	mappedUserObject[username] = &newUser
}

// GetUserStateInfo get *user.RealUser corresponding to username
func GetUserStateInfo(username string) *user.RealUser {
	_, ok := mappedUserObject[username]
	fmt.Printf("[Storage]get user: %s is ok=%v\n", username, ok)
	if !ok {
		ensureUser(username)
	}
	return mappedUserObject[username]
}

var mappedDeviceObject = make(map[string]*device.Device)

func ensureDevice(deviceID string) {
	common.Println("[Storage] make new device for", deviceID)
	mdb, err := common.Mongo()
	if common.PrintError(err) {
		fmt.Println("  At ensureDevice()")
		return
	}
	var tmp map[string]interface{}
	err = mdb.DB("CUSmartFarm").C("devices").Find(bson.M{
		"id": deviceID,
	}).One(&tmp)
	if err != nil {
		return
	}
	dev := device.Device{}
	dev.FromMap(tmp)
	mappedDeviceObject[deviceID] = &dev
}

// GetDevice get device object
func GetDevice(deviceID string) (dev *device.Device, err error) {
	_, ok := mappedDeviceObject[deviceID]
	if !ok { // then make the new device
		ensureDevice(deviceID)
	}
	return mappedDeviceObject[deviceID], nil
}
