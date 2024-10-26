package service

import (
	"Service-schema/domain"
	"Service-schema/persistence"
	"Service-schema/service/dto"
	"errors"
	"fmt"
)

type IProductService interface {
	Add(createProductRequestDto dto.CreateProductRequestDto) error
	Delete(productId int64) error
	UpdatePrice(updateProductRequestDto dto.UpdateProductRequestDto) error
	GetById(productId int64) (domain.Product, error)
	GetAllProducts() []domain.Product
	GetAllProductsByStoreName(storeName string) []domain.Product
}

type ProductService struct {
	productRepository persistence.IProductRepository
}

func NewProductService(productRepository persistence.IProductRepository) IProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (productService *ProductService) Add(createProductRequestDto dto.CreateProductRequestDto) error {
	validationErr := validateCreateProductRequestDto(createProductRequestDto)
	if validationErr != nil {
		return validationErr
	}

	return productService.productRepository.Add(domain.Product{
		Name:     createProductRequestDto.Name,
		Price:    createProductRequestDto.Price,
		Discount: createProductRequestDto.Discount,
		Store:    createProductRequestDto.Store,
	})
}

func (productService *ProductService) Delete(productId int64) error {
	return productService.productRepository.DeleteById(productId)
}

func (productService *ProductService) UpdatePrice(updateProductRequestDto dto.UpdateProductRequestDto) error {

	validationErr := validateUpdateProductRequestDto(updateProductRequestDto)

	if validationErr != nil {
		return validationErr
	}

	return productService.productRepository.UpdatePrice(updateProductRequestDto.Id, updateProductRequestDto.Price)
}

func (productService *ProductService) GetById(productId int64) (domain.Product, error) {
	return productService.productRepository.GetById(productId)
}

func (productService *ProductService) GetAllProducts() []domain.Product {
	return productService.productRepository.GetAllProducts()
}

func (productService *ProductService) GetAllProductsByStoreName(storeName string) []domain.Product {
	return productService.productRepository.GetAllProductsByStoreName(storeName)
}

func validateCreateProductRequestDto(createProductRequestDto dto.CreateProductRequestDto) error {
	if createProductRequestDto.Discount > float32(50.0) {
		return errors.New(fmt.Sprintf("Discount must be less than 50 percent"))
	}

	if createProductRequestDto.Store == "" {
		return errors.New(fmt.Sprintf("Store must be specified"))
	}

	if createProductRequestDto.Name == "" {
		return errors.New(fmt.Sprintf("Name must be specified"))
	}

	if createProductRequestDto.Price < float32(10.0) {
		return errors.New(fmt.Sprintf("Price must be greater than 10"))
	}

	return nil
}

func validateUpdateProductRequestDto(updateProductRequestDto dto.UpdateProductRequestDto) error {
	if updateProductRequestDto.Id == 0 {
		return errors.New(fmt.Sprintf("Id must be specified"))
	}

	if updateProductRequestDto.Price < float32(10.0) {
		return errors.New(fmt.Sprintf("Price must be greater than 10"))
	}

	return nil
}
