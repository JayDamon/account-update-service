package config

import (
	"fmt"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerplaid"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	HostPort        string
	ApplicationName string
	KeyCloakConfig  *moneymakergocloak.Configuration
	Plaid           *moneymakerplaid.Configuration
	Rabbit          *Configuration
}

func GetConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	hostPort := getOrDefault("HOST_PORT", "3000")
	applicationName := getOrDefault("APPLICATION_NAME", "")

	config := &Config{
		HostPort:        hostPort,
		ApplicationName: applicationName,
		//KeyCloakConfig: moneymakergocloak.NewConfiguration(),
		//Plaid:          moneymakerplaid.NewConfiguration(),
		Rabbit: NewConfiguration(),
	}

	return config
}

func getOrDefault(envVar string, defaultVal string) string {
	val := os.Getenv(envVar)
	if val == "" {
		return defaultVal
	}
	return val
}
