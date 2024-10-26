package service

import (
	"Service-schema/domain"
	"Service-schema/persistence"
	"Service-schema/service"
	"Service-schema/service/dto"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var productService service.IProductService

func TestMain(m *testing.M) {
	var initializedProducts = []domain.Product{
		{
			Id:       1,
			Name:     `360Hz 24" Monitor`,
			Price:    2000.0,
			Discount: 12.0,
			Store:    "BENQ",
		},
		{
			Id:       2,
			Name:     `EC-2B Mouse`,
			Price:    1200.0,
			Discount: 10.0,
			Store:    "Zowie",
		},
		{
			Id:       3,
			Name:     `RTX 5090`,
			Price:    10000.0,
			Discount: 20.0,
			Store:    "Nvidia",
		},
		{
			Id:       4,
			Name:     `Iphone 17`,
			Price:    3000.0,
			Discount: 0.0,
			Store:    "Apple",
		},
	}
	fakeProductRepository := NewFakeProductRepository(initializedProducts)
	productService = service.NewProductService(persistence.IProductRepository(fakeProductRepository))

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_ShouldGetAllProducts(t *testing.T) {
	t.Run("ShouldGetAllProducts", func(t *testing.T) {
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
	})
}

func Test_ShouldGetAllProductsByStoreName(t *testing.T) {
	t.Run("ShouldGetAllProductsByStoreName", func(t *testing.T) {
		actualProducts := productService.GetAllProductsByStoreName("BENQ")
		assert.Equal(t, 1, len(actualProducts))
	})
}

func Test_ShouldGetProductById(t *testing.T) {
	expectedProduct := domain.Product{
		Id:       4,
		Name:     `Iphone 17`,
		Price:    3000.0,
		Discount: 0.0,
		Store:    "Apple",
	}
	t.Run("ShouldGetProductById", func(t *testing.T) {
		actualProduct, _ := productService.GetById(4)
		assert.Equal(t, expectedProduct, actualProduct)
		_, err := productService.GetById(10)
		assert.Equal(t, "Product not found", err.Error())
	})
}

func Test_WhenNoValidationErrorOccurred_ShouldAddProduct(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldAddProduct", func(t *testing.T) {
		productRequest := dto.CreateProductRequestDto{
			Name:     "Pencil",
			Price:    200.0,
			Discount: 0.0,
			Store:    "Amazon",
		}
		_ = productService.Add(productRequest)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 5, len(actualProducts))
	})
	_ = productService.Delete(5)
}

func Test_WhenNameFieldEmpty_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenNameFieldEmpty_ShouldNotAddProduct", func(t *testing.T) {
		productRequest := dto.CreateProductRequestDto{
			Name:     "",
			Price:    200.0,
			Discount: 0.0,
			Store:    "Amazon",
		}

		err := productService.Add(productRequest)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, "Name must be specified", err.Error())
	})
}

func Test_WhenStoreFieldEmpty_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenStoreFieldEmpty_ShouldNotAddProduct", func(t *testing.T) {
		productRequest := dto.CreateProductRequestDto{
			Name:     "Pencil",
			Price:    200.0,
			Discount: 0.0,
			Store:    "",
		}

		err := productService.Add(productRequest)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, "Store must be specified", err.Error())
	})

}

func Test_WhenDiscountFieldHigherThan50_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenDiscountFieldHigherThan50_ShouldNotAddProduct", func(t *testing.T) {
		productRequest := dto.CreateProductRequestDto{
			Name:     "Pencil",
			Price:    200.0,
			Discount: 51.0,
			Store:    "Amazon",
		}

		err := productService.Add(productRequest)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, "Discount must be less than 50 percent", err.Error())
	})
}

func Test_WhenPriceFieldLessThan10_ShouldNotAddProduct(t *testing.T) {
	t.Run("WhenPriceFieldLessThan10_ShouldNotAddProduct", func(t *testing.T) {
		productRequest := dto.CreateProductRequestDto{
			Name:     "Pencil",
			Price:    2.0,
			Discount: 0.0,
			Store:    "Amazon",
		}

		err := productService.Add(productRequest)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, "Price must be greater than 10", err.Error())
	})
}

func Test_WhenNoValidationErrorOccurred_ShouldUpdateProductPrice(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldUpdateProductPrice", func(t *testing.T) {
		updatedProductRequest := dto.UpdateProductRequestDto{
			Id:    4,
			Price: 4000.0,
		}
		_ = productService.UpdatePrice(updatedProductRequest)
		actualProduct, _ := productService.GetById(4)
		assert.Equal(t, updatedProductRequest.Price, actualProduct.Price)
	})
}

func Test_WhenIdFieldEmpty_ShouldNotUpdateProductPrice(t *testing.T) {
	t.Run("WhenIdFieldEmpty_ShouldNotUpdateProductPrice", func(t *testing.T) {
		updatedProductRequest := dto.UpdateProductRequestDto{
			Price: 4000.0,
		}
		err := productService.UpdatePrice(updatedProductRequest)
		assert.Equal(t, "Id must be specified", err.Error())
	})
}

func Test_WhenPriceFieldLessThan10_ShouldNotUpdateProductPrice(t *testing.T) {
	t.Run("WhenPriceFieldLessThan10_ShouldNotUpdateProductPrice", func(t *testing.T) {
		updatedProductRequest := dto.UpdateProductRequestDto{
			Id:    4,
			Price: 4.0,
		}
		err := productService.UpdatePrice(updatedProductRequest)
		actualProduct, _ := productService.GetById(4)
		assert.Equal(t, "Price must be greater than 10", err.Error())
		assert.Equal(t, float32(3000.0), actualProduct.Price)
	})
}

func Test_WhenGivenCorrectProductId_ShouldDeleteProduct(t *testing.T) {
	t.Run("WhenGivenCorrectProductId_ShouldDeleteProduct", func(t *testing.T) {
		_ = productService.Delete(4)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 3, len(actualProducts))
	})
}

func Test_WhenGivenWrongProductId_ShouldNotDeleteProduct(t *testing.T) {
	t.Run("WhenGivenCorrectProductId_ShouldDeleteProduct", func(t *testing.T) {
		err := productService.Delete(10)
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, "Product not found", err.Error())
	})
}
