package repository

import (
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/example/model"
)

type Response map[string]interface{}

type DeviceError struct {
}

func (d *DeviceError) Process(ctx context.Context, code int, msg string) {
	logger.Log.Error(code, msg)
}

func GetDevice(c *gin.Context) {
	device_id := c.Param("device_id")
	resp := &model.GetDeviceResponse{}
	err := connector.MakeGetRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/iot-03/devices/%s", device_id)),
		connector.WithResp(resp),
		connector.WithErrProc(1102, &DeviceError{}))
	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		c.Abort()
		return
	}
	c.JSON(200, resp)
}

func PostDeviceCmd(c *gin.Context) {
	device_id := c.Param("device_id")
	body, _ := io.ReadAll(c.Request.Body)
	resp := &model.PostDeviceCmdResponse{}
	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/iot-03/devices/%s/commands", device_id)),
		connector.WithPayload(body),
		connector.WithResp(resp))
	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		c.Abort()
		return
	}
	c.JSON(200, resp)
}
