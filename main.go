package main

import (
	"fmt"
	"github.com/p4xx07/order-service/configuration"
	"github.com/p4xx07/order-service/deps"
	"github.com/p4xx07/order-service/internal/log"
)

func main() {
	c, err := configuration.GetEnvConfig()
	if err != nil {
		fmt.Println("failed to load configuration")
		panic(err)
	}

	zapLogger := log.NewLogger(c.LogLevel)
	app, err := deps.InjectApp(c, zapLogger)
	if err != nil {
		zapLogger.Fatal(err)
	}

	err = app.Routes().Listen("0.0.0.0:8080")
	zapLogger.Error(err)
}
