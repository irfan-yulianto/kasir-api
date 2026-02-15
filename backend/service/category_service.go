package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type CategoryService interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(id int, category *model.Category) error
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Create(category *model.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) Update(id int, category *model.Category) error {
	return s.repo.Update(id, category)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
