package models

import "time"

type ProductImage struct {
	ID        string
	ProductID string
	ImageURL  string
	Position  int
	CreatedAt time.Time
}
