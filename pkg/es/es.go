package es

import (
	"context"
	"errors"

	Conf "github.com/annaworks/chatservice/pkg/conf"

	"github.com/olivere/elastic"
	"go.uber.org/zap"
)

type es struct {
	host string
	client *elastic.Client
}

func NewElasticSearch(conf Conf.Conf) *es {
	return &es{
		host: conf.ES_HOST,
		client: nil,
	}
}

func (e *es) SetupES() error {
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(e.host),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return err
	}

	e.client = elasticClient

	return nil
}


func (e es) CreateDb(name, schema string) error {
	ctx := context.Background()

	res, err := e.client.CreateIndex(name).BodyString(schema).Do(ctx)
	if err != nil {
		return err
	}

	if !res.Acknowledged {
		return errors.New("Registering was not acknowledged")
	}

	return nil
}

func (e es) DbExists(name string) (bool, error) {
	ctx := context.Background()

	exists, err := e.client.IndexExists(name).Do(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}


