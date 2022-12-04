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
	_, err := c.conn.Exec("INSERT INTO cars (model, brand_id, city, year, price, description) VALUES ($1, $2, $3, $4, $5, $6)",
		car.Model, car.BrandID, car.City, car.Year, car.Price, car.Description)
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

func (c CarsRepository) Sort(ctx context.Context, sortType string) ([]*models.Car, error) {
	sortedCars := make([]*models.Car, 0)
	if sortType == "model-asc" {
		if err := c.conn.Select(&sortedCars, "SELECT * FROM cars ORDER BY model;"); err != nil {
			return nil, err
		}
	}
	if sortType == "price-asc" {
		if err := c.conn.Select(&sortedCars, "SELECT * FROM cars ORDER BY price;"); err != nil {
			return nil, err
		}
	}
	return sortedCars, nil
}

func (c CarsRepository) FilterByCity(ctx context.Context, filter string) ([]*models.Car, error) {
	filteredCars := make([]*models.Car, 0)
	if err := c.conn.Select(&filteredCars, "SELECT * FROM cars WHERE city ILIKE $1", ""+filter+""); err != nil {
		return nil, err
	}
	return filteredCars, nil
}

func (c CarsRepository) AddToFav(ctx context.Context, filter *models.CarFilter) error {
	favouriteCar := new(models.Car)
	basicQuery := "SELECT * FROM cars WHERE id = $1"

	if filter.CarId != nil {
		if err := c.conn.Get(favouriteCar, basicQuery, filter.CarId); err != nil {
			return err
		}
	}

	_, err := c.conn.Exec("INSERT INTO favourites(car_id) VALUES ($1)", favouriteCar.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c CarsRepository) DeleteFromFav(ctx context.Context, filter *models.CarFilter) error {
	favouriteCar := new(models.Car)
	basicQuery := "SELECT * FROM cars WHERE id = $1"

	if filter.CarId != nil {
		if err := c.conn.Get(favouriteCar, basicQuery, filter.CarId); err != nil {
			return err
		}
	}

	_, err := c.conn.Exec("DELETE FROM favourites WHERE id = $1", favouriteCar.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c CarsRepository) ShowFav(ctx context.Context) ([]*models.Car, error) {
	favouriteCars := make([]*models.Car, 0)
	err := c.conn.Select(&favouriteCars, "select cars.id, cars.model, cars.brand_id, cars.city, cars.year, cars.description from cars, favourites where cars.id =  favourites.car_id")
	if err != nil {
		return nil, err
	}
	return favouriteCars, nil
}
