package config

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var ErrReadInConfig error = errors.New("error while reading configuration from yaml file")
var ErrDecodeConfig error = errors.New("unable to decode configuration")

// Device represents each Tuya device with a Name and ID
type Device struct {
	Name string `mapstructure:"name"`
	ID   string `mapstructure:"id"`
}

type Ambilight struct {
	// In milliseconds
	RefreshRate int `mapstructure:"ambilight.refreshRate"`
}

// Config represents the application's configuration, including Tuya devices
type Config struct {
	AccessID              string    `mapstructure:"accessID"`
	AccessKey             string    `mapstructure:"accessKey"`
	AppName               string    `mapstructure:"appName"`
	DebugMode             bool      `mapstructure:"debugMode"`
	TuyaRegisteredDevices []Device  `mapstructure:"tuya.devices"`
	Ambilight             Ambilight `mapstructure:"ambilight"`
}

// LoadConfig loads configuration from a file or environment variables
func LoadConfig() (*Config, error) {
	var cfg Config

	// Set the configuration file name and type
	viper.SetConfigName("config")  // Config file name (without extension)
	viper.SetConfigType("yaml")    // Config file type (e.g., yaml, json, toml, etc.)
	viper.AddConfigPath("config/") // Add an additional path to search for the config file

	// Environment variables override file settings if they exist
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("Failed to read configuration file")
		return nil, ErrReadInConfig
	}
	log.Info().Msg("Configuration file loaded successfully")

	// Map file/environment values to the Config struct
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode configuration into struct")
		return &Config{}, ErrDecodeConfig
	}

	// Manually retrieve the registered devices
	rawDevices := viper.Get("tuya.devices")
	if rawDevices != nil {
		deviceMaps := rawDevices.([]interface{}) // Type assertion to []interface{}
		for _, deviceMap := range deviceMaps {
			device := deviceMap.(map[string]interface{})                                              // Assert to map[string]interface{}
			name := device["name"].(string)                                                           // Extract Name
			id := device["id"].(string)                                                               // Extract ID
			cfg.TuyaRegisteredDevices = append(cfg.TuyaRegisteredDevices, Device{Name: name, ID: id}) // Add to slice
		}
	}

	log.Info().
		Str("AccessID", cfg.AccessID).
		Str("AppName", cfg.AppName).
		Int("DeviceCount", len(cfg.TuyaRegisteredDevices)).
		Bool("DebugMode", cfg.DebugMode).
		Msg("Configuration loaded and mapped to struct")

	return &cfg, nil
}
