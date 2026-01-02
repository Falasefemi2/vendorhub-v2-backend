package dto

import (
	"errors"
)

// CreateProductRequest - Request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"required,min=1,max=1000"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

// Validate validates the create product request
func (r *CreateProductRequest) Validate() error {
	if r.Name == "" {
		return errors.New("product name is required")
	}
	if len(r.Name) > 255 {
		return errors.New("product name must be less than 255 characters")
	}
	if r.Description == "" {
		return errors.New("product description is required")
	}
	if len(r.Description) > 1000 {
		return errors.New("product description must be less than 1000 characters")
	}
	if r.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	return nil
}

// UpdateProductRequest - Request body for updating a product
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	IsActive    *bool    `json:"is_active"`
}

// Validate validates the update product request
func (r *UpdateProductRequest) Validate() error {
	if r.Name != nil && len(*r.Name) > 255 {
		return errors.New("product name must be less than 255 characters")
	}
	if r.Description != nil && len(*r.Description) > 1000 {
		return errors.New("product description must be less than 1000 characters")
	}
	if r.Price != nil && *r.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	return nil
}

// ProductResponse - Response body for product
type ProductResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
