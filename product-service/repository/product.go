package repository

import (
	"context"

	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IProductRepository interface {
	Create(ctx context.Context, data model.Product) error
	FindListProductQttyData(ctx context.Context, productIds []string) ([]model.ProductQttyData, error)
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

func (r *productRepository) FindListProductQttyData(ctx context.Context, productIds []string) ([]model.ProductQttyData, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: productIds}}}}
	projection := bson.D{{Key: "_d", Value: 1}, {Key: "qtty", Value: 1}}
	list := []model.ProductQttyData{}

	cursor, err := r.client.
		Database(constant.DATABASE_NAME_PRODUCT_SERVICE).
		Collection(constant.COLLECTION_NAME_PRODUCTS).
		Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return list, err
	}

	for cursor.Next(ctx) {
		item := model.ProductQttyData{}
		if err := cursor.Decode(&item); err != nil {
			return list, err
		}

		list = append(list, item)
	}

	if err := cursor.Err(); err != nil {
		return list, err
	}

	return list, nil
}
