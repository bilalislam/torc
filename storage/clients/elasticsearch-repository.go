package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/bilalislam/torc/storage"
	"github.com/bilalislam/torc/storage/models"
	"github.com/olivere/elastic/v7"
)

type IElasticsearchRepository interface {
	storage.IRepository
}

type ElasticsearchRepository struct {
	Client *elastic.Client
}

func (er *ElasticsearchRepository) GetById(ctx context.Context, id string, model models.IModel) error {
	return nil
}

func (er *ElasticsearchRepository) Save(ctx context.Context, model models.IModel) error {
	if er.Client == nil {
		return errors.New("elasticsearch client or connection nil. Please initialize first")
	}
	_, err := er.Client.Index().
		Index("syslog-2021.11.16").
		BodyJson(model).
		Do(ctx)

	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			panic(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			panic(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			panic(fmt.Sprintf("Connection problem: %v", err))
		default:
			return err
		}
	}

	return nil
}

func (er *ElasticsearchRepository) Update(ctx context.Context, model models.IModel) (int64, error) {
	return 0, nil
}

func (er *ElasticsearchRepository) Delete(ctx context.Context, id string) (int64, error) {
	return 0, nil
}

func (er *ElasticsearchRepository) Connect(ctx context.Context, dbContext interface{}) error {
	return nil
}
