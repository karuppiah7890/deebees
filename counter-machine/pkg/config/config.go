package config

import (
	"fmt"
	"os"
	"strconv"
)

// All configuration is through environment variables

const PORT_ENV_VAR = "PORT"
const DEFAULT_PORT = 8080

type Config struct {
	port int
}

func NewConfigFromEnvVars() (*Config, error) {
	port, err := getPort()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting application port: %v", err)
	}

	return &Config{
		port: port,
	}, nil
}

// Get application port number. Default is 8080
func getPort() (int, error) {
	portEnvStr, ok := os.LookupEnv(PORT_ENV_VAR)
	if !ok {
		return DEFAULT_PORT, nil
	}

	port, err := strconv.Atoi(portEnvStr)
	if err != nil {
		return 0, fmt.Errorf("%s environment variable value (%s) cannot be converted to integer", PORT_ENV_VAR, portEnvStr)
	}

	return port, nil
}

func (c *Config) GetPort() int {
	return c.port
}
