package conf

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Conf struct {
	// api
	API_PORT string `envconfig:"API_PORT" default:"3000"`

	// slack
	SLACK_TOKEN              string `envconfig:"SLACK_TOKEN"`
}

func NewConf(logger *zap.Logger) Conf {
	c := &Conf{}

	err := envconfig.Process("", c)
	if err != nil {
		logger.Fatal("Envconfig could not get env", zap.Error(err))
	}

	logger.Info("Envconfig success")

	return *c
}
