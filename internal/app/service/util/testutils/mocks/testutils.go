package testutils

import (
	"log"
	"testing"
	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/service/logger"
	"uniswapper/internal/config"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func SetupTest(t *testing.T, envPath string) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	constants.Config, err = config.LoadConfig(envPath)
	assert.NoError(t, err)
	logger.InitLogger()

}
