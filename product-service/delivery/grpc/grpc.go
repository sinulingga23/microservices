package grpc

import (
	pbProduct "github.com/sinulingga23/microservices/product-service/proto-generated/product"
	"github.com/sinulingga23/microservices/product-service/service/rpc"
	"google.golang.org/grpc"
)

func InitRegistrationServer(productRpc *rpc.ProductRpc) *grpc.Server {
	grpcServer := grpc.NewServer()

	pbProduct.RegisterProductServer(grpcServer, productRpc)

	return grpcServer
}
