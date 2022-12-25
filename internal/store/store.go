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
	Users() UsersRepository
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
	AllOfUser(ctx context.Context, userId int) ([]*models.Car, error)
	ByID(ctx context.Context, id int) (*models.Car, error)
	Update(ctx context.Context, car *models.Car) error
	Delete(ctx context.Context, id int) error
	Sort(ctx context.Context, sortType string) ([]*models.Car, error)
	FilterByCity(ctx context.Context, filter string) ([]*models.Car, error)
	AddToFav(ctx context.Context, filter *models.CarFilter) error
	ShowFav(ctx context.Context) ([]*models.Car, error)
	DeleteFromFav(ctx context.Context, filter *models.CarFilter) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context) ([]*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}
