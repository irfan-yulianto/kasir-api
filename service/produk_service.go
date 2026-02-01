package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type ProdukService interface {
	GetAll() ([]model.Produk, error)
	GetAllWithCategory() ([]model.Produk, error)
	GetByID(id int) (*model.Produk, error)
	GetByIDWithCategory(id int) (*model.Produk, error)
	GetByCategoryID(categoryID int) ([]model.Produk, error)
	Create(produk *model.Produk) error
	Update(id int, produk *model.Produk) error
	Delete(id int) error
}

type produkService struct {
	repo repository.ProdukRepository
}

func NewProdukService(repo repository.ProdukRepository) ProdukService {
	return &produkService{repo: repo}
}

func (s *produkService) GetAll() ([]model.Produk, error) {
	return s.repo.GetAll()
}

func (s *produkService) GetAllWithCategory() ([]model.Produk, error) {
	return s.repo.GetAllWithCategory()
}

func (s *produkService) GetByID(id int) (*model.Produk, error) {
	return s.repo.GetByID(id)
}

func (s *produkService) GetByIDWithCategory(id int) (*model.Produk, error) {
	return s.repo.GetByIDWithCategory(id)
}

func (s *produkService) GetByCategoryID(categoryID int) ([]model.Produk, error) {
	return s.repo.GetByCategoryID(categoryID)
}

func (s *produkService) Create(produk *model.Produk) error {
	return s.repo.Create(produk)
}

func (s *produkService) Update(id int, produk *model.Produk) error {
	return s.repo.Update(id, produk)
}

func (s *produkService) Delete(id int) error {
	return s.repo.Delete(id)
}
