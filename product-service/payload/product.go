package payload

import "github.com/sinulingga23/microservices/product-service/utils"

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

	return nil
}
