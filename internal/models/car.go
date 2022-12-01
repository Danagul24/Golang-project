package models

type (
	Car struct {
		ID          int    `json:"id" db:"id"`
		Model       string `json:"model" db:"model"`
		BrandID     int    `json:"brand_id" db:"brand_id"`
		City        string `json:"city" db:"city"`
		Year        int    `json:"year" db:"year"`
		Price       int    `json:"price" db:"price"`
		Description string `json:"description" db:"description"`
	}

	CarFilter struct {
		Query *string `json:"query"`
	}
)
