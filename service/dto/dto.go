package dto

type CreateProductRequestDto struct {
	Name     string
	Price    float32
	Discount float32
	Store    string
}

type UpdateProductRequestDto struct {
	Id    int64
	Price float32
}
