package repository

import (
	"context"

	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type IProductRepository interface {
	Create(ctx context.Context, data model.Product) error
}

type productRepository struct {
	client *mongo.Client
}

func NewProductRepository(client *mongo.Client) IProductRepository {
	return &productRepository{client: client}
}

func (r *productRepository) Create(ctx context.Context, data model.Product) error {
	_, err := r.client.
		Database(constant.DATABASE_NAME_PRODUCT_SERVICE).
		Collection(constant.COLLECTION_NAME_PRODUCTS).
		InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
