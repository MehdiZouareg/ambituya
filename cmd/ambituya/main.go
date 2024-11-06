package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/httplib"
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

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode is enabled")
	}

	// Init Tuya Connector
	log.Info().Msg("Initializing Tuya connector")
	connector.InitWithOptions(
		env.WithApiHost(httplib.URL_EU),
		env.WithMsgHost(httplib.MSG_EU),
		env.WithAccessID(config.AccessID),
		env.WithAccessKey(config.AccessKey),
		env.WithAppName(config.AppName),
		env.WithDebugMode(config.DebugMode),
	)

	// Declare flag variables
	var deviceID, command, switchOn string

	// Flag definition
	flag.StringVar(&deviceID, "device-id", "", "Specify the device ID")
	flag.StringVar(&command, "command", "", "Specify the command")
	flag.StringVar(&switchOn, "switch", "", "Specify the command")
	flag.Parse()

	// Check whether flags are used, if not we just to send one command to our devices and shutdown process
	if deviceID == "" && command == "" {
		log.Info().Msg("No command-line flags provided, running in normal mode")

		if len(config.TuyaRegisteredDevices) == 0 {
			log.Fatal().Msg("no devices configured in config file.")
		}

		go func() {
			log.Info().Msg("starting system tray...")
			systray.Systray(config)
		}()

		go func() {
			log.Info().Msg("Starting message listener...")
			messaging.Listener()
		}()

		go func() {
			log.Info().Msg("Starting Ambilight effect...")
			effect.Ambilight(config)
		}()

		waitForExitSignal()
	} else {
		if switchOn != "" {
			st, err := strconv.ParseBool(switchOn)
			if err != nil {
				os.Exit(0)
			}
			tuya.Switch(deviceID, st)
		}

		tuya.CastCommandToID(deviceID, command)

		os.Exit(0)
	}
}

func waitForExitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Info().Str("Signal", sig.String()).Msg("Received termination signal, shutting down gracefully")
}
