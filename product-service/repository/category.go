package repository

import (
	"context"

	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryRepository interface {
	Create(ctx context.Context, data model.Category) error
}

type categoryRepository struct {
	client *mongo.Client
}

func NewCategoryRepository(client *mongo.Client) ICategoryRepository {
	return &categoryRepository{client: client}
}

func (r *categoryRepository) Create(ctx context.Context, data model.Category) error {
	_, err := r.client.
		Database(constant.DATABASE_NAME_PRODUCT_SERVICE).
		Collection(constant.COLLECTION_NAME_CATAGORIES).
		InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
