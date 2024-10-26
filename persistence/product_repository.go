package persistence

import (
	"Service-schema/domain"
	"Service-schema/persistence/common"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStoreName(storeName string) []domain.Product
	GetById(productId int64) (domain.Product, error)
	Add(product domain.Product) error
	DeleteById(productId int64) error
	UpdatePrice(productId int64, price float32) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func (productRepository *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	selectQuery := "SELECT * FROM products"
	productRows, err := productRepository.dbPool.Query(ctx, selectQuery)

	if err != nil {
		log.Error("Error occurred getting all products %v", err)
		return []domain.Product{}
	}

	return extractsAllProducts(productRows)
}

func (productRepository *ProductRepository) GetAllProductsByStoreName(storeName string) []domain.Product {
	ctx := context.Background()
	selectProductsByStoreNameQuery := `SELECT * FROM products WHERE store=$1`
	productRows, err := productRepository.dbPool.Query(ctx, selectProductsByStoreNameQuery, storeName)

	if err != nil {
		log.Error("Error occurred getting all products %v", err)
		return []domain.Product{}
	}

	return extractsAllProducts(productRows)
}

func (productRepository *ProductRepository) Add(product domain.Product) error {
	ctx := context.Background()
	insertSQL := `insert into products (name,price,discount,store) values ($1,$2,$3,$4)`
	newProduct, err := productRepository.dbPool.Exec(ctx, insertSQL, product.Name, product.Price, product.Discount, product.Store)
	if err != nil {
		log.Error("Error occurred inserting product %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Added product %v", newProduct))
	return nil
}

func (productRepository *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()
	selectQuery := `SELECT * FROM products WHERE id = $1`
	productRow := productRepository.dbPool.QueryRow(ctx, selectQuery, productId)
	return extractProduct(productId, productRow)
}

func (productRepository *ProductRepository) DeleteById(productId int64) error {
	ctx := context.Background()
	_, productGetErr := productRepository.GetById(productId)

	if productGetErr != nil {
		return errors.New(fmt.Sprintf("Product not found"))
	}

	deleteSQL := `DELETE FROM products WHERE id = $1`
	_, err := productRepository.dbPool.Exec(ctx, deleteSQL, productId)
	if err != nil {
		log.Error("Error occurred deleting product %v", err)
		return errors.New(fmt.Sprintf("Error occurred deleting product with id %d", productId))
	}

	log.Info(fmt.Sprintf("Deleted product %d", productId))
	return nil
}

func (productRepository *ProductRepository) UpdatePrice(productId int64, price float32) error {
	ctx := context.Background()
	_, productGetErr := productRepository.GetById(productId)
	if productGetErr != nil {
		return errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}

	updateSQL := `UPDATE products SET price = $2 WHERE id = $1`
	_, err := productRepository.dbPool.Exec(ctx, updateSQL, productId, price)
	if err != nil {
		log.Error("Error occurred updating product %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Updated product %d", productId))
	return nil
}

func extractsAllProducts(productRows pgx.Rows) []domain.Product {
	var products []domain.Product
	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	for productRows.Next() {
		productRows.Scan(&id, &name, &price, &discount, &store)
		product := domain.Product{Id: id, Name: name, Price: price, Discount: discount, Store: store}
		products = append(products, product)
	}

	return products
}

func extractProduct(productId int64, productRow pgx.Row) (domain.Product, error) {

	var id int64
	var name string
	var price float32
	var discount float32
	var store string
	scanErr := productRow.Scan(&id, &name, &price, &discount, &store)

	if scanErr != nil && scanErr.Error() == common.NOT_FOUND {
		return domain.Product{}, errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	if scanErr != nil {
		return domain.Product{}, errors.New(fmt.Sprintf("Error occurred when scanned product with id %d", productId))
	}

	product := domain.Product{Id: id, Name: name, Price: price, Discount: discount, Store: store}
	return product, nil
}
