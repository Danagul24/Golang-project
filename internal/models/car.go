package models

import "math/big"

type Car struct {
	ID          int     `json:"id" db:"id"`
	Model       string  `json:"model" db:"model"`
	Brand       string  `json:"brand" db:"brand"`
	City        string  `json:"city" db:"city"`
	Year        int     `json:"year" db:"year"`
	Price       big.Int `json:"price" db:"price"`
	Description string  `json:"description" db:"description"`
}
