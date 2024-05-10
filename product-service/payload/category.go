package payload

import (
	"github.com/sinulingga23/microservices/product-service/utils"
)

type AddCategoryRequest struct {
	Name string `json:"name"`
}

func (p *AddCategoryRequest) Validate() error {
	if p.Name == "" {
		return utils.ErrNameCantEmpty
	}

	return nil
}
