package rpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sinulingga23/microservices/product-service/constant"
	pbBase "github.com/sinulingga23/microservices/product-service/proto-generated/base"
	pbProduct "github.com/sinulingga23/microservices/product-service/proto-generated/product"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRpc struct {
	pbProduct.UnimplementedProductServer

	productRepository repository.IProductRepository
}

func NewProductRpc(productRepository repository.IProductRepository) *ProductRpc {
	return &ProductRpc{productRepository: productRepository}
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

	data := in.Data

	uniqueIds := []string{}
	visitedId := map[string]int{}
	for _, item := range data {
		_, ok := visitedId[item.ProductId]
		if !ok {
			uniqueIds = append(uniqueIds, item.ProductId)
		}

		visitedId[item.ProductId] = 1
	}

	listProductQttyData, err := r.productRepository.FindListProductQttyData(ctx, uniqueIds)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		response.ResponseCode = constant.RPC_PRODUCT_NOT_EXISTS
		response.ResponseDesc = utils.ErrProductNotExists.Error()
		return response, nil
	}
	if err != nil {
		response.ResponseCode = constant.RPC_CODE_FAILED
		response.ResponseDesc = utils.ErrFailed.Error()
		return response, nil
	}

	if len(listProductQttyData) == 0 {
		response.ResponseCode = constant.RPC_PRODUCT_NOT_EXISTS
		response.ResponseDesc = utils.ErrProductNotExists.Error()
		return response, nil
	}

	count := 0
	for _, item := range listProductQttyData {
		_, ok := visitedId[item.Id]
		if ok {
			count += 1
		}
	}

	if count != len(uniqueIds) {
		response.ResponseCode = constant.RPC_PRODUCT_NOT_EXISTS
		response.ResponseDesc = utils.ErrProductNotExists.Error()
		return response, nil
	}

	// TODO: Dedeuct qtty on db

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
