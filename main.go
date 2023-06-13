package main

import (
	"context"
	"fmt"
	"uniswapper/internal/app/api/server"
	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/db"
	"uniswapper/internal/app/service/logger"
	"uniswapper/internal/config"
)

func main() {
	var err error
	// Returns a struct with values from env variables
	constants.Config, err = config.LoadConfig(".env")
	if err != nil {
		panic(err.Error())
	}
	// Creates an empty context that can be passed around
	ctx := context.Background()

	// Initialize the logger
	logger.InitLogger()
	log := logger.Logger(ctx)

	// Creates db connection & applies migrations
	dbConn, err := db.Init(ctx)
	if err != nil {
		log.Fatalf("DB connection failed with error: %v", err)
	}
	dbConnection := db.New(dbConn)

	r := server.Init(ctx, dbConnection)
	if err := r.Run(fmt.Sprintf("%s:%s", constants.Config.HTTPServerConfig.HTTPSERVER_LISTEN, constants.Config.HTTPServerConfig.HTTPSERVER_PORT)); err != nil {
		log.Fatal("Server not able to startup with error: ", err)
	}
}
