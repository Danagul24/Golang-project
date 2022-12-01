package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"project/internal/models"
	"project/internal/store"
)

func (db *DB) Brands() store.BrandsRepository {
	if db.brands == nil {
		db.brands = newBrandsRepository(db.conn)
	}
	return db.brands
}

type BrandsRepository struct {
	conn *sqlx.DB
}

func newBrandsRepository(conn *sqlx.DB) store.BrandsRepository {
	return &BrandsRepository{conn: conn}
}

func (c BrandsRepository) Create(ctx context.Context, brand *models.Brand) error {
	_, err := c.conn.Exec("INSERT INTO brands(name) VALUES ($1)", brand.Name)
	if err != nil {
		return err
	}
	return nil
}

func (c BrandsRepository) All(ctx context.Context, filter *models.BrandFilter) ([]*models.Brand, error) {
	brands := make([]*models.Brand, 0)
	basicQuery := "SELECT * FROM brands"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := c.conn.Select(&brands, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}
		return brands, nil
	}

	if err := c.conn.Select(&brands, basicQuery); err != nil {
		return nil, err
	}
	return brands, nil
}

func (c BrandsRepository) ByID(ctx context.Context, id int) (*models.Brand, error) {
	brand := new(models.Brand)
	if err := c.conn.Get(brand, "SELECT id, name FROM brands WHERE id = $1", id); err != nil {
		return nil, err
	}
	return brand, nil
}

func (c BrandsRepository) Update(ctx context.Context, brand *models.Brand) error {
	_, err := c.conn.Exec("UPDATE brands SET name = $1 WHERE id = $2", brand.Name, brand.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c BrandsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM brands WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
