package main

import (
	"fmt"
	"log"

	"github.com/annaworks/chatservice/pkg/api"
	Conf "github.com/annaworks/chatservice/pkg/conf"

	"go.uber.org/zap"
)

func main() {
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout"}
	logger, err := c.Build()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not init zap logger: %v", err))
	}
	defer logger.Sync()

	// env
	conf := Conf.NewConf(logger.Named("conf_logger"))

	// api
	api := api.NewApi(logger.Named("api_logger"), conf)
	api.Init()
	api.Serve()
}
