package rpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/model"
	pbBase "github.com/sinulingga23/microservices/product-service/proto-generated/base"
	pbProduct "github.com/sinulingga23/microservices/product-service/proto-generated/product"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ProductRpc struct {
	pbProduct.UnimplementedProductServer

	client            *mongo.Client
	productRepository repository.IProductRepository
}

func NewProductRpc(
	productRepository repository.IProductRepository,
	client *mongo.Client) *ProductRpc {
	return &ProductRpc{
		productRepository: productRepository,
		client:            client,
	}
}

func (r *ProductRpc) HandleDeductQtty(ctx context.Context, in *pbProduct.DeductQttyRequest) (*pbBase.BaseResponse, error) {
	response := &pbBase.BaseResponse{
		ResponseCode: "00",
		ResponseDesc: "Success",
	}

	if rc, rd := r.validateDeductQttyRequest(in); rc != "" && rd != "" {
		response.ResponseCode = rc
		response.ResponseDesc = rd
		return response, nil
	}

	wc := writeconcern.Majority()
	txOptions := options.Transaction().SetWriteConcern(wc)

	txSession, err := r.client.StartSession()
	if err != nil {
		log.Printf("Failed start session of tx: %v", err)
		response.ResponseCode = constant.RPC_CODE_FAILED
		response.ResponseDesc = constant.RPC_CODE_FAILED
		return response, nil
	}
	defer txSession.EndSession(context.TODO())

	result, err := txSession.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		productCollection := r.client.
			Database(constant.DATABASE_NAME_PRODUCT_SERVICE).
			Collection(constant.COLLECTION_NAME_PRODUCTS)

		ordersProductsCollection := r.client.
			Database(constant.COLLECTION_NAME_PRODUCTS).
			Collection(constant.COLLECTION_NAME_ORDERS_PRODUCTS)

		lenData := len(in.Data)
		for i := 0; i < lenData; i++ {
			item := (in.Data)[i]

			singleResult := productCollection.
				FindOne(ctx,
					bson.D{{Key: "_id", Value: item.ProductId}},
					options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}, {Key: "qtty", Value: 1}}))

			currentData := model.ProductQttyData{}
			if err := singleResult.Decode(&currentData); err != nil {
				log.Printf("error when try decoce data: %v", err)
				return err, nil
			}

			if err := singleResult.Err(); err != nil {
				log.Printf("there last error: %v", err)
				return err, nil
			}

			if currentData.Qtty < int(item.Qtty) {
				return utils.ErrInsufficientQtty, nil
			}
			currentData.Qtty -= int(item.Qtty)

			_, err = productCollection.UpdateByID(ctx, bson.D{{Key: "_id", Value: item.ProductId}}, bson.D{{Key: "$set", Value: bson.D{{Key: "qtty", Value: currentData.Qtty}}}})
			if err != nil {
				log.Printf("error when try update product by id: %v", err)
				return utils.ErrFailed, nil
			}

			orderProduct := model.OrderProduct{
				OrderId:   in.OrderId,
				ProductId: item.ProductId,
				Qtty:      int(item.Qtty),
			}
			_, err = ordersProductsCollection.InsertOne(ctx, orderProduct)
			if err != nil {
				log.Printf("error when try insert new data orders_product: %v", err)
				return utils.ErrFailed, nil
			}
		}

		return nil, nil
	}, txOptions)
	if err != nil {
		log.Printf("Failed when do tx: %v", err)
		response.ResponseCode = constant.RPC_CODE_FAILED
		response.ResponseDesc = constant.RPC_CODE_FAILED
		return response, nil
	}

	log.Println(result)

	return &pbBase.BaseResponse{}, nil
}

func (r *ProductRpc) GetListProductByIds(ctx context.Context, in *pbProduct.GetListProductByIdsRequest) (*pbBase.BaseResponse, error) {
	response := &pbBase.BaseResponse{
		ResponseCode: "00",
		ResponseDesc: "Success",
	}

	// TODO: Implementation

	return response, nil
}

func (r *ProductRpc) validateDeductQttyRequest(in *pbProduct.DeductQttyRequest) (rc string, rd string) {
	data := in.Data
	if data == nil {
		return constant.RPC_CODE_DATA_EMPTY, utils.ErrDataEmpty.Error()
	}

	for _, item := range data {
		_, err := uuid.Parse(item.ProductId)
		isNotValidRequest := err != nil || item.Qtty <= 0
		if isNotValidRequest {
			return constant.RPC_CODE_INVALID_REQUEST, utils.ErrInvalidRequest.Error()
		}
	}

	return "", ""
}
