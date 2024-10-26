package request

import "Service-schema/service/dto"

type CreateProductRequest struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Discount float32 `json:"discount"`
	Store    string  `json:"store"`
}

type UpdateProductPriceRequest struct {
	Id    int64   `json:"id"`
	Price float32 `json:"price"`
}

func (createProductRequest CreateProductRequest) ToDto() dto.CreateProductRequestDto {
	return dto.CreateProductRequestDto{
		Name:     createProductRequest.Name,
		Price:    createProductRequest.Price,
		Discount: createProductRequest.Discount,
		Store:    createProductRequest.Store,
	}
}

func (updateProductRequest UpdateProductPriceRequest) ToDto() dto.UpdateProductRequestDto {
	return dto.UpdateProductRequestDto{
		Id:    updateProductRequest.Id,
		Price: updateProductRequest.Price,
	}
}
