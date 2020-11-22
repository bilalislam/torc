package storage

import (
	"context"
	"github.com/bilalislam/torc/storage/models"
)

//todo redis do not need context. refactor it.
type IRepository interface {
	GetById(ctx context.Context, id string, model models.IModel) error
	Save(ctx context.Context, model models.IModel) error
	Update(ctx context.Context, model models.IModel) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Connect(ctx context.Context, dbContext interface{}) error
}
