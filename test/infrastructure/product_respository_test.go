package infrastructure

import (
	"Service-schema/core/postgresql"
	"Service-schema/domain"
	"Service-schema/persistence"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var productRepository persistence.IProductRepository
var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		UserName:              "postgres",
		Password:              "password",
		DbName:                "product_service",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	})

	productRepository = persistence.NewProductRepository(dbPool)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()
	setup(ctx, dbPool)
	expected := []domain.Product{
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
	t.Run("TestGetAllProducts", func(t *testing.T) {
		actualProducts := productRepository.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, expected, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestGetAllProductsByStoreName(t *testing.T) {
	ctx := context.Background()
	setup(ctx, dbPool)
	expected := []domain.Product{
		{
			Id:       3,
			Name:     `RTX 5090`,
			Price:    10000.0,
			Discount: 20.0,
			Store:    "Nvidia",
		},
	}
	t.Run("TestGetAllProductsByStoreName", func(t *testing.T) {
		actualProducts := productRepository.GetAllProductsByStoreName("Nvidia")
		assert.Equal(t, 1, len(actualProducts))
		assert.Equal(t, expected, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	expected := []domain.Product{
		{
			Id:       1,
			Name:     "Samsung Galaxy A72",
			Price:    5000.0,
			Discount: 0.0,
			Store:    "Samsung",
		},
	}
	t.Run("TestAddProduct", func(t *testing.T) {
		newProduct := domain.Product{
			Name:     "Samsung Galaxy A72",
			Price:    5000.0,
			Discount: 0.0,
			Store:    "Samsung",
		}

		_ = productRepository.Add(newProduct)

		products := productRepository.GetAllProducts()
		assert.Equal(t, 1, len(products))
		assert.Equal(t, expected, products)
	})
	clear(ctx, dbPool)
}

func TestGetById(t *testing.T) {
	ctx := context.Background()
	setup(ctx, dbPool)
	expectedProduct := domain.Product{
		Id:       3,
		Name:     `RTX 5090`,
		Price:    10000.0,
		Discount: 20.0,
		Store:    "Nvidia",
	}
	t.Run("TestGetById", func(t *testing.T) {
		actualProduct, _ := productRepository.GetById(3)
		_, err := productRepository.GetById(5)
		assert.Equal(t, expectedProduct, actualProduct)
		assert.Equal(t, "Product not found with id 5", err.Error())
	})
	clear(ctx, dbPool)
}

func TestDeleteById(t *testing.T) {
	ctx := context.Background()
	setup(ctx, dbPool)
	t.Run("TestDeleteById", func(t *testing.T) {
		_ = productRepository.DeleteById(3)
		products := productRepository.GetAllProducts()
		assert.Equal(t, 3, len(products))
	})
	clear(ctx, dbPool)
}

func TestUpdateProductPrice(t *testing.T) {
	ctx := context.Background()
	setup(ctx, dbPool)
	t.Run("TestUpdateProductPrice", func(t *testing.T) {
		_ = productRepository.UpdatePrice(3, 12000.0)
		err := productRepository.UpdatePrice(5, 5000)
		actualProduct, _ := productRepository.GetById(3)
		assert.Equal(t, float32(12000.0), actualProduct.Price)
		assert.Equal(t, "Product not found with id 5", err.Error())
	})
	clear(ctx, dbPool)
}

func setup(ctx context.Context, dbPool *pgxpool.Pool) {
	TestDataInitialize(ctx, dbPool)
}

func clear(ctx context.Context, dbPool *pgxpool.Pool) {
	TruncateTestData(ctx, dbPool)
}
