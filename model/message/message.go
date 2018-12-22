package message

import "encoding/json"

// APICall represent payload for API Calls via WS
type APICall struct {
	EndPoint string                 `json:"endPoint" binding:"required"` // addDevice, removeDevice, setDevice, pollDevice, listDevice...
	Token    string                 `json:"token" binding:"required"`
	Payload  map[string]interface{} `json:"payload" binding:"required"` // json data depend on command
}

// RelayState represents state of a relay (On/Off) and it's detail
type RelayState struct { // use when set relay mode
	Mode   string      `json:"mode" binding:"required"` // ON OFF AUTO SCHED ...
	Detail interface{} `json:"detail"`                  // detail depending on mode
}

// AddDeviceMessage is payload format for addDevice API
type AddDeviceMessage struct {
	DeviceID     string `json:"deviceID" binding:"required"`
	DeviceSecret string `json:"deviceSecret" binding:"required"`
}

// RemoveDeviceMessage is payload format for removeDevice API
type RemoveDeviceMessage struct {
	DeviceID string `json:"deviceID" binding:"required"`
}

// DeviceCommandMessage is payload format for setDevice API
type DeviceCommandMessage struct {
	DeviceID string     `json:"deviceID" binding:"required"`
	RelayID  string     `json:"relayID" binding:"required"`
	State    RelayState `json:"state" binding:"required"`
}

// FromMap is "constructor" for converting map[string]interface{} to AddDeviceMessage  return error if can't convert
func (message *AddDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

// FromMap is "constructor" for converting map[string]interface{} to RemoveDeviceMessage  return error if can't convert
func (message *RemoveDeviceMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}

// FromMap is "constructor" for converting map[string]interface{} to DeviceCommandMessage  return error if can't convert
func (message *DeviceCommandMessage) FromMap(val map[string]interface{}) error {
	str, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, message)
	if err != nil {
		return err
	}
	return nil
}