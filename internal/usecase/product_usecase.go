package usecase

import "github.com/supakorn141/golang-erp/internal/domain"

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.productRepo.FindAll()
}

func (u *productUsecase) GetByID(id uint) (*domain.Product, error) {
	return u.productRepo.FindByID(id)
}

func (u *productUsecase) Create(product *domain.Product) error {
	return u.productRepo.Create(product)
}

func (u *productUsecase) Update(id uint, input *domain.Product) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	product.Name = input.Name
	product.SKU = input.SKU
	product.Price = input.Price
	product.Stock = input.Stock
	product.Category = input.Category
	return u.productRepo.Update(product)
}

func (u *productUsecase) Delete(id uint) error {
	return u.productRepo.Delete(id)
}
