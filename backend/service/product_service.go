package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type ProductService interface {
	GetAll() ([]model.Product, error)
	GetAllWithCategory() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	GetByIDWithCategory(id int) (*model.Product, error)
	GetByCategoryID(categoryID int) ([]model.Product, error)
	Create(product *model.Product) error
	Update(id int, product *model.Product) error
	Delete(id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetAllWithCategory() ([]model.Product, error) {
	return s.repo.GetAllWithCategory()
}

func (s *productService) GetByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) GetByIDWithCategory(id int) (*model.Product, error) {
	return s.repo.GetByIDWithCategory(id)
}

func (s *productService) GetByCategoryID(categoryID int) ([]model.Product, error) {
	return s.repo.GetByCategoryID(categoryID)
}

func (s *productService) Create(product *model.Product) error {
	return s.repo.Create(product)
}

func (s *productService) Update(id int, product *model.Product) error {
	return s.repo.Update(id, product)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
