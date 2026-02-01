package repository

import (
	"database/sql"
	"kasir-api/model"
)

type ProdukRepository interface {
	GetAll() ([]model.Produk, error)
	GetAllWithCategory() ([]model.Produk, error)
	GetByID(id int) (*model.Produk, error)
	GetByIDWithCategory(id int) (*model.Produk, error)
	GetByCategoryID(categoryID int) ([]model.Produk, error)
	Create(produk *model.Produk) error
	Update(id int, produk *model.Produk) error
	Delete(id int) error
}

type produkRepository struct {
	db *sql.DB
}

func NewProdukRepository(db *sql.DB) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) GetAll() ([]model.Produk, error) {
	query := "SELECT id, nama, harga, stok, category_id FROM produk ORDER BY id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Produk
	for rows.Next() {
		var p model.Produk
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GetAllWithCategory menggunakan SQL JOIN untuk mendapatkan produk beserta kategorinya
func (r *produkRepository) GetAllWithCategory() ([]model.Produk, error) {
	query := `
		SELECT 
			p.id, p.nama, p.harga, p.stok, p.category_id,
			c.id, c.name, c.description
		FROM produk p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Produk
	for rows.Next() {
		var p model.Produk
		var catID, catIDFromJoin sql.NullInt64
		var catName, catDesc sql.NullString

		if err := rows.Scan(
			&p.ID, &p.Nama, &p.Harga, &p.Stok, &catID,
			&catIDFromJoin, &catName, &catDesc,
		); err != nil {
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

func (r *produkRepository) GetByID(id int) (*model.Produk, error) {
	query := "SELECT id, nama, harga, stok, category_id FROM produk WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var p model.Produk
	var catID sql.NullInt64
	if err := row.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &catID); err != nil {
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

// GetByIDWithCategory menggunakan SQL JOIN untuk mendapatkan produk beserta kategorinya
func (r *produkRepository) GetByIDWithCategory(id int) (*model.Produk, error) {
	query := `
		SELECT 
			p.id, p.nama, p.harga, p.stok, p.category_id,
			c.id, c.name, c.description
		FROM produk p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	row := r.db.QueryRow(query, id)

	var p model.Produk
	var catID, catIDFromJoin sql.NullInt64
	var catName, catDesc sql.NullString

	if err := row.Scan(
		&p.ID, &p.Nama, &p.Harga, &p.Stok, &catID,
		&catIDFromJoin, &catName, &catDesc,
	); err != nil {
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

// GetByCategoryID mendapatkan semua produk berdasarkan category_id
func (r *produkRepository) GetByCategoryID(categoryID int) ([]model.Produk, error) {
	query := `
		SELECT 
			p.id, p.nama, p.harga, p.stok, p.category_id,
			c.id, c.name, c.description
		FROM produk p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.category_id = $1
		ORDER BY p.id
	`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Produk
	for rows.Next() {
		var p model.Produk
		var catID, catIDFromJoin sql.NullInt64
		var catName, catDesc sql.NullString

		if err := rows.Scan(
			&p.ID, &p.Nama, &p.Harga, &p.Stok, &catID,
			&catIDFromJoin, &catName, &catDesc,
		); err != nil {
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

func (r *produkRepository) Create(produk *model.Produk) error {
	query := "INSERT INTO produk (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	return r.db.QueryRow(query, produk.Nama, produk.Harga, produk.Stok, produk.CategoryID).Scan(&produk.ID)
}

func (r *produkRepository) Update(id int, produk *model.Produk) error {
	query := "UPDATE produk SET nama = $1, harga = $2, stok = $3, category_id = $4 WHERE id = $5"
	result, err := r.db.Exec(query, produk.Nama, produk.Harga, produk.Stok, produk.CategoryID, id)
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

	produk.ID = id
	return nil
}

func (r *produkRepository) Delete(id int) error {
	query := "DELETE FROM produk WHERE id = $1"
	result, err := r.db.Exec(query, id)
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
