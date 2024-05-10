package payload

import (
	"github.com/google/uuid"
	"github.com/sinulingga23/microservices/product-service/utils"
)

type AddProductRequest struct {
	Name        string   `json:"name"`
	Qtty        int      `json:"qtty"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	CategoryIds []string `json:"categoryIds"`
}

func (p *AddProductRequest) Validate() error {
	if p.Name == "" {
		return utils.ErrNameCantEmpty
	}

	if p.Qtty <= 0 {
		return utils.ErrQttyInvalid
	}
	if p.Price <= 0 {
		return utils.ErrPriceInvalid
	}

	for i := 0; i < len(p.CategoryIds); i++ {
		id := p.CategoryIds[i]
		_, err := uuid.Parse(id)
		if err != nil {
			return utils.ErrCategoryNotExists
		}
	}

	return nil
}
