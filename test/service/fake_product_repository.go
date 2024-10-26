package service

import (
	"Service-schema/domain"
	"Service-schema/persistence"
	"errors"
	"fmt"
)

type FakeProductRepository struct {
	products []domain.Product
}

var currentIdValue int64

func NewFakeProductRepository(initializeProducts []domain.Product) persistence.IProductRepository {
	currentIdValue = int64(len(initializeProducts)) + 1
	return &FakeProductRepository{
		products: initializeProducts,
	}
}

func (fakeRepository *FakeProductRepository) GetAllProducts() []domain.Product {
	return fakeRepository.products
}

func (fakeRepository *FakeProductRepository) GetAllProductsByStoreName(storeName string) []domain.Product {
	var productsWithStoreName []domain.Product
	for index, product := range fakeRepository.products {
		if product.Store == storeName {
			productsWithStoreName = append(productsWithStoreName, fakeRepository.products[index])
		}
	}
	return productsWithStoreName
}

func (fakeRepository *FakeProductRepository) Add(product domain.Product) error {
	fakeRepository.products = append(fakeRepository.products, domain.Product{
		Id:       currentIdValue,
		Name:     product.Name,
		Price:    product.Price,
		Discount: product.Discount,
		Store:    product.Store,
	})
	currentIdValue++
	return nil
}

func (fakeRepository *FakeProductRepository) GetById(productId int64) (domain.Product, error) {
	for index, product := range fakeRepository.products {
		if product.Id == productId {
			return fakeRepository.products[index], nil
		}
	}
	return domain.Product{}, errors.New(fmt.Sprintf("Product not found"))
}

func (fakeRepository *FakeProductRepository) DeleteById(productId int64) error {

	for index, product := range fakeRepository.products {
		if product.Id == productId {
			fakeRepository.products = append(fakeRepository.products[:index-1], fakeRepository.products[index:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Product not found"))
}

func (fakeRepository *FakeProductRepository) UpdatePrice(productId int64, price float32) error {
	for index, product := range fakeRepository.products {
		if product.Id == productId {
			fakeRepository.products[index].Price = price
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Product not found"))
}
