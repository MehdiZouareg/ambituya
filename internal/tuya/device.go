package tuya

import (
	"context"
	"errors"
	"fmt"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/logger"
)

var v1 InstructionSet
var v2 InstructionSet

type InstructionSet struct {
	Name         string
	Instructions map[string]Instruction
}

type Instruction struct {
	Name   string
	MinVal int
	MaxVal int
	*Instruction
}

type StatusKey struct {
	Code  string
	Value interface{}
}

type GetDevicesResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	T       int64       `json:"t"`
}

type GetDeviceResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	T       int64       `json:"t"`
}

func InitInstructionsSets() {
	v1 = InstructionSet{
		Name: "v1",
		Instructions: map[string]Instruction{
			"color":      {Name: "colour_data", MinVal: 0, MaxVal: 255},
			"brightness": {Name: "bright_value", MinVal: 10, MaxVal: 1000},
		},
	}

	v2 = InstructionSet{
		Name: "v2",
		Instructions: map[string]Instruction{
			"color":      {Name: "colour_data_v2", MinVal: 0, MaxVal: 255},
			"brightness": {Name: "bright_value_v2", MinVal: 10, MaxVal: 1000}},
	}
}

var RegisteredDevices []Device

type GetAllDeviceResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	T       int64       `json:"t"`
}

// Should retrieve all devices registered on a Tuya Project
func GetAllDevicesInProject() (*[]Device, error) {
	resp := &GetAllDeviceResponse{}

	err := connector.MakeGetRequest(
		context.Background(),
		connector.WithAPIUri("/v2.0/cloud/thing/device"),
		connector.WithResp(resp))
	if err != nil {
		// fmt.Println("err:", err.Error())
		return nil, err
	}

	return nil, errors.New("couldnt get project devices")
}

func (d *Device) GetDeviceStatus(code string) (interface{}, error) {
	resp := &GetDeviceResponse{}

	err := connector.MakeGetRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s", d.ID)),
		connector.WithResp(resp))
	if err != nil {
		return "", err
	}

	if resp.Result != nil {
		for k, v := range resp.Result.(map[string]interface{}) {
			if k == "status" {
				for _, j := range v.([]interface{}) {
					Value := j.(map[string]interface{})["value"]

					return Value, nil
				}
			}
		}
	}

	return "", errors.New("couldnt get device status")
}

func GetDevicesList(ids []string) ([]Device, error) {
	resp := &GetDevicesResponse{}

	for _, i := range ids {
		err := connector.MakeGetRequest(
			context.Background(),
			connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%v", i)),
			connector.WithResp(resp))
		if err != nil {
			return nil, err
		}
	}

	return resp.Result.([]Device), nil
}

func (d *Device) Switch() error {
	status, err := d.GetDeviceStatus("switch_led")
	if err != nil {
		return err.(error)
	}
	command := fmt.Sprintf(`{
		"commands": [
		{
			"code": "switch_led",
			"value": %v
	  	}
	]}`, !status.(bool))

	err = connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", d.ID)),
		connector.WithPayload([]byte(command)))
	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		return err
	}

	return nil
}
