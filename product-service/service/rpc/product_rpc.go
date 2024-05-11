package rpc

import (
	pbProduct "github.com/sinulingga23/microservices/product-service/proto-generated/product"
)

type productRpc struct {
	pbProduct.UnimplementedProductServer
}
