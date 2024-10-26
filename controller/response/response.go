package response

import "Service-schema/domain"

type ErrorResponse struct {
	ErrorDescription string `json:"error_description"`
}

type ProductResponse struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

func ToProductResponse(product domain.Product) ProductResponse {
	return ProductResponse{
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
		Store:    product.Store,
	}
}

func ToProductResponseList(products []domain.Product) []ProductResponse {
	var productResponses = []ProductResponse{}

	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}

	return productResponses
}
