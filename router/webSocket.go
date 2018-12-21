package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"github.com/rod41732/cu-smart-farm-backend/model"
	User "github.com/rod41732/cu-smart-farm-backend/model/user"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

var wsUpgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     wsCheckOrigin,
}

func wsCheckOrigin(r *http.Request) bool {
	return true // todo origin check
}

// WebSocket : WS request handler
func WebSocket(c *gin.Context) {
	// tmp, _ := c.Get("username")
	// user := tmp.(middleware.User)
	user := middleware.User{Username: "rod41732"}
	wsRouter(c.Writer, c.Request, user.Username)
}

func wsRouter(w http.ResponseWriter, r *http.Request, username string) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	user := User.New(username, "random", conn)

	if common.PrintError(err) {
		return
	}
	for { // loop for command
		t, msg, err := conn.ReadMessage()
		if err != nil {

			break
		}

		var payload model.APICall
		err = json.Unmarshal(msg, payload)

		var success bool
		var errmsg string
		if err != nil {
			success, errmsg = false, "No Command Specified"
		} else {
			switch payload.EndPoint {
			case "addDevice":
				user.AddDevice(payload.Payload)
			case "removeDevice":
				user.RemoveDevice(payload.Payload)
			case "setDevice":
			}
		}

		resp, err := json.Marshal(bson.M{"success": success, "errmsg": errmsg})
		conn.WriteMessage(t, resp)
		time.Sleep(time.Millisecond * 10)
	}

}
