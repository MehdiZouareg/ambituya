package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config représente la configuration de l'application
type Config struct {
	AccessID  string
	AccessKey string
	AppName   string
	DebugMode bool
}

// LoadConfig charge la configuration à partir d'un fichier ou de variables d'environnement
func LoadConfig() (*Config, error) {
	var cfg Config

	// Définir le fichier de configuration
	viper.SetConfigName("config")  // Nom du fichier de config (sans l'extension)
	viper.SetConfigType("yaml")    // Type du fichier de config (ici JSON, mais ça pourrait être yaml, toml, etc.)
	viper.AddConfigPath("config/") // Ajoute aussi un autre chemin pour plus de flexibilité

	// Lire les variables d'environnement en priorité si elles existent
	viper.AutomaticEnv()

	// Lire le fichier de config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Erreur lors du chargement du fichier de config: %v\n", err)
	}

	// Associer les valeurs du fichier/variables d'environnement à la structure Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("impossible de décoder la configuration: %w", err)
	}

	return &cfg, nil
}
