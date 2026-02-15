package repository

import (
	"database/sql"
	"kasir-api/model"
)

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	GetAllWithCategory() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	GetByIDWithCategory(id int) (*model.Product, error)
	GetByCategoryID(categoryID int) ([]model.Product, error)
	Create(product *model.Product) error
	Update(id int, product *model.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock, category_id FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) GetAllWithCategory() ([]model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
			   c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var catID, catIDFromJoin sql.NullInt64
		var catName, catDesc sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID,
			&catIDFromJoin, &catName, &catDesc); err != nil {
			return nil, err
		}

		if catID.Valid {
			categoryID := int(catID.Int64)
			p.CategoryID = &categoryID
		}

		if catIDFromJoin.Valid {
			p.Category = &model.Category{
				ID:          int(catIDFromJoin.Int64),
				Name:        catName.String,
				Description: catDesc.String,
			}
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product
	var catID sql.NullInt64

	err := r.db.QueryRow("SELECT id, name, price, stock, category_id FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if catID.Valid {
		categoryID := int(catID.Int64)
		p.CategoryID = &categoryID
	}
	return &p, nil
}

func (r *productRepository) GetByIDWithCategory(id int) (*model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
			   c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1`

	var p model.Product
	var catID, catIDFromJoin sql.NullInt64
	var catName, catDesc sql.NullString

	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID,
		&catIDFromJoin, &catName, &catDesc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if catID.Valid {
		categoryID := int(catID.Int64)
		p.CategoryID = &categoryID
	}

	if catIDFromJoin.Valid {
		p.Category = &model.Category{
			ID:          int(catIDFromJoin.Int64),
			Name:        catName.String,
			Description: catDesc.String,
		}
	}
	return &p, nil
}

func (r *productRepository) GetByCategoryID(categoryID int) ([]model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
			   c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.category_id = $1
		ORDER BY p.id`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var catID, catIDFromJoin sql.NullInt64
		var catName, catDesc sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID,
			&catIDFromJoin, &catName, &catDesc); err != nil {
			return nil, err
		}

		if catID.Valid {
			cID := int(catID.Int64)
			p.CategoryID = &cID
		}

		if catIDFromJoin.Valid {
			p.Category = &model.Category{
				ID:          int(catIDFromJoin.Int64),
				Name:        catName.String,
				Description: catDesc.String,
			}
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Price, product.Stock, product.CategoryID,
	).Scan(&product.ID)
}

func (r *productRepository) Update(id int, product *model.Product) error {
	result, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
		product.Name, product.Price, product.Stock, product.CategoryID, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	product.ID = id
	return nil
}

func (r *productRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
