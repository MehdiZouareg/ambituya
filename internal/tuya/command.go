package tuya

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tuya/tuya-connector-go/connector"
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
		log.Error().Err(err).Msg("got error while sending switch command to device")
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
		log.Error().Err(err).Msg("got error while casting command to device")
		return err
	}

	return nil
}
