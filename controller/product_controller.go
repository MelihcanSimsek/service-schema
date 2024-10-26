package controller

import (
	"Service-schema/controller/request"
	"Service-schema/controller/response"
	"Service-schema/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProductController struct {
	productService service.IProductService
}

func NewProductController(productService service.IProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (productController *ProductController) RegisterRoutes(e *echo.Echo) {

	e.GET("/api/v1/products/:id", productController.GetProductById)
	e.GET("/api/v1/products/", productController.GetAllProducts)
	e.POST("/api/v1/products/", productController.Add)
	e.PUT("/api/v1/products/", productController.UpdatePrice)
	e.DELETE("/api/v1/products/:id", productController.DeleteProductById)
}

func (productController *ProductController) GetProductById(c echo.Context) error {
	id := c.Param("id")
	productId, convertErr := strconv.Atoi(id)

	if convertErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: convertErr.Error()})
	}

	product, err := productController.productService.GetById(int64(productId))

	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{ErrorDescription: err.Error()})
	}

	return c.JSON(http.StatusOK, response.ToProductResponse(product))
}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	store := c.QueryParam("store")

	if len(store) == 0 {
		products := productController.productService.GetAllProducts()
		return c.JSON(http.StatusOK, response.ToProductResponseList(products))
	}

	productsWithStoreName := productController.productService.GetAllProductsByStoreName(store)
	return c.JSON(http.StatusOK, response.ToProductResponseList(productsWithStoreName))
}

func (productController *ProductController) Add(c echo.Context) error {
	var createProductRequest request.CreateProductRequest
	bindErr := c.Bind(&createProductRequest)
	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: bindErr.Error()})
	}

	err := productController.productService.Add(createProductRequest.ToDto())

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{ErrorDescription: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func (productController *ProductController) UpdatePrice(c echo.Context) error {
	var updateProductPriceRequest request.UpdateProductPriceRequest
	bindErr := c.Bind(&updateProductPriceRequest)

	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: bindErr.Error()})
	}

	err := productController.productService.UpdatePrice(updateProductPriceRequest.ToDto())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{ErrorDescription: err.Error()})
	}

	return c.NoContent(http.StatusAccepted)
}

func (productController *ProductController) DeleteProductById(c echo.Context) error {
	id := c.Param("id")
	productId, converErr := strconv.Atoi(id)
	if converErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{ErrorDescription: converErr.Error()})
	}

	err := productController.productService.Delete(int64(productId))

	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{ErrorDescription: err.Error()})
	}

	return c.NoContent(http.StatusAccepted)
}
