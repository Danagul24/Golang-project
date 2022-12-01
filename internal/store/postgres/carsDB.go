package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"project/internal/models"
	"project/internal/store"
)

func (db *DB) Cars() store.CarsRepository {
	if db.cars == nil {
		db.cars = newCarsRepository(db.conn)
	}
	return db.cars
}

type CarsRepository struct {
	conn *sqlx.DB
}

func newCarsRepository(conn *sqlx.DB) store.CarsRepository {
	return &CarsRepository{conn: conn}
}

func (c CarsRepository) Create(ctx context.Context, car *models.Car) error {
	_, err := c.conn.Exec("INSERT INTO cars(model, brand_id, city, year, price, description)"+
		" VALUES ($1, $2, $3, &4, $5, &6)", car.Model, car.BrandID, car.City, car.Year, car.Price, car.Description)
	if err != nil {
		return err
	}
	return nil
}

func (c CarsRepository) All(ctx context.Context, filter *models.CarFilter) ([]*models.Car, error) {
	cars := make([]*models.Car, 0)
	basicQuery := "SELECT * FROM cars"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE model ILIKE $1", basicQuery)

		if err := c.conn.Select(&cars, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}
		return cars, nil
	}

	if err := c.conn.Select(&cars, basicQuery); err != nil {
		return nil, err
	}
	return cars, nil
}

func (c CarsRepository) ByID(ctx context.Context, id int) (*models.Car, error) {
	car := new(models.Car)
	if err := c.conn.Get(car, "SELECT * FROM cars WHERE id = $1", id); err != nil {
		return nil, err
	}
	return car, nil
}

func (c CarsRepository) Update(ctx context.Context, car *models.Car) error {
	_, err := c.conn.Exec("UPDATE cars SET city = $1 WHERE id = $2", car.City, car.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c CarsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
