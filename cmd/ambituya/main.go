package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/env/extension"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/example/messaging"

	"github.com/MehdiZouareg/ambituya/config"
	"github.com/MehdiZouareg/ambituya/internal/effect"
	"github.com/MehdiZouareg/ambituya/internal/systray"
	"github.com/MehdiZouareg/ambituya/internal/tuya"
)

func main() {
	// Start with loading configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("couldnt load initialization config")
	}

	// Init Tuya Connector
	connector.InitWithOptions(
		env.WithApiHost(httplib.URL_EU),
		env.WithMsgHost(httplib.MSG_EU),
		env.WithAccessID(config.AccessID),
		env.WithAccessKey(config.AccessKey),
		env.WithAppName(config.AppName),
		env.WithDebugMode(false),
	)

	// Declare flag variables
	var deviceID string
	var command string
	var switchOn string

	// Flag definition
	flag.StringVar(&deviceID, "device-id", "", "Specify the device ID")
	flag.StringVar(&command, "command", "", "Specify the command")
	flag.StringVar(&switchOn, "switch", "", "Specify the command")

	flag.Parse()

	// Check whether flags are used
	if deviceID == "" && command == "" {
		// fmt.Println("No flags provided.")

		logger.Log.SetLevel(logger.ERROR)

		tuya.RegisteredDevices = []tuya.Device{
			{
				Name: "LEDS",
				ID:   "545000508cce4ee25ac7",
			}, {
				Name: "Salon",
				ID:   "06325004e868e74c5f14",
			},
		}

		go systray.Systray()

		go messaging.Listener()

		go effect.Ambilight(tuya.RegisteredDevices)

		waitSignal()
	} else {
		if switchOn != "" {
			st, err := strconv.ParseBool(switchOn)
			if err != nil {
				os.Exit(0)
			}
			Switch(deviceID, st)
		}

		CastCommandToID(deviceID, command)

		os.Exit(0)
	}
}

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
	fmt.Printf("\n\nSending %v to %v", command, id)

	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", id)),
		connector.WithPayload([]byte(command)))
	if err != nil {
		logger.Log.Errorf("err:%s", err.(error).Error())
		return err.(error)
	}

	return nil
}

func waitSignal() {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		select {
		case c := <-quitCh:
			extension.GetMessage(constant.TUYA_MESSAGE).Stop()
			logger.Log.Infof("receive sig:%v, shutdown the http server...", c.String())
			return
		}
	}
}
