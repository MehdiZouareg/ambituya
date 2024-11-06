package tuya

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/logger"
)

func Switch(id string, state bool) error {
	command := fmt.Sprintf(`{
		"commands": [
		{
			"code": "switch_led",
			"value": %v
	  	}
	]}`, state)

	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", id)),
		connector.WithPayload([]byte(command)))
	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		return err
	}

	return nil
}

func CastCommandToID(id, command string) error {
	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", id)),
		connector.WithPayload([]byte(command)))
	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		return err
	}

	return nil
}
