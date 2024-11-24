package config

import (
	"fmt"
	"github.com/jaydamon/moneymakergocloak"
	"github.com/jaydamon/moneymakerplaid"
	"github.com/jaydamon/moneymakerrabbit"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	HostPort        string
	ApplicationName string
	ConfigureCors   bool
	KeyCloakConfig  *moneymakergocloak.Configuration
	Plaid           *moneymakerplaid.Configuration
	Rabbit          *moneymakerrabbit.Configuration
}

func GetConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	hostPort := getOrDefault("HOST_PORT", "3000")
	applicationName := getOrDefault("APPLICATION_NAME", "")

	configureCors := getOrDefaultBool("CONFIGURE_CORS", true)
	config := &Config{
		HostPort:        hostPort,
		ApplicationName: applicationName,
		ConfigureCors:   configureCors,
		KeyCloakConfig:  moneymakergocloak.NewConfiguration(),
		Plaid:           moneymakerplaid.NewConfiguration(),
		Rabbit:          moneymakerrabbit.NewConfiguration(),
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

func getOrDefaultBool(envVar string, defaultVal bool) bool {
	val := os.Getenv(envVar)
	var returnVal = defaultVal
	if val == "true" {
		returnVal = true
	} else if val == "false" {
		returnVal = false
	}

	return returnVal
}
