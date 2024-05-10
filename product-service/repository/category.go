package repository

import (
	"context"

	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICategoryRepository interface {
	Create(ctx context.Context, data model.Category) error
	GetIdsByIds(ctx context.Context, ids []string) ([]string, error)
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

func (r *categoryRepository) GetIdsByIds(ctx context.Context, ids []string) ([]string, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	projection := bson.D{{Key: "_id", Value: 1}}

	listId := make([]string, 0)

	cursor, err := r.client.
		Database(constant.DATABASE_NAME_PRODUCT_SERVICE).
		Collection(constant.COLLECTION_NAME_CATAGORIES).
		Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return listId, err
	}

	for cursor.Next(ctx) {
		mapId := map[string]string{}
		if err := cursor.Decode(&mapId); err != nil {
			return listId, err
		}

		listId = append(listId, mapId["_id"])
	}

	if err := cursor.Err(); err != nil {
		return listId, err
	}

	return listId, nil
}
