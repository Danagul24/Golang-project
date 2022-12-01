package store

import (
	"context"
	"project/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error
	Brands() BrandsRepository
	Cars() CarsRepository
}

type BrandsRepository interface {
	Create(ctx context.Context, brand *models.Brand) error
	All(ctx context.Context, filter *models.BrandFilter) ([]*models.Brand, error)
	ByID(ctx context.Context, id int) (*models.Brand, error)
	Update(ctx context.Context, brand *models.Brand) error
	Delete(ctx context.Context, id int) error
}

type CarsRepository interface {
	Create(ctx context.Context, car *models.Car) error
	All(ctx context.Context, filter *models.CarFilter) ([]*models.Car, error)
	ByID(ctx context.Context, id int) (*models.Car, error)
	Update(ctx context.Context, car *models.Car) error
	Delete(ctx context.Context, id int) error
}
