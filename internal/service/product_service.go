package service

import (
	"context"
	"fmt"
	"time"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/models"
	"github.com/falasefemi2/vendorhub/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(ctx context.Context, vendorID string, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	product := &models.Product{
		UserID:      vendorID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsActive:    true,
	}

	createdProduct, err := ps.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return mapProductToResponse(createdProduct), nil
}

// GetProduct retrieves a single product by ID
func (ps *ProductService) GetProduct(ctx context.Context, productID string) (*dto.ProductResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	product, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return mapProductToResponse(product), nil
}

// GetUserProducts retrieves all products for a user
func (ps *ProductService) GetUserProducts(ctx context.Context, userID string) ([]*dto.ProductResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetProductsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// UpdateProduct updates an existing product
func (ps *ProductService) UpdateProduct(ctx context.Context, productID string, vendorID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if productID == "" || vendorID == "" {
		return nil, fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Verify product exists and belongs to vendor
	existingProduct, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if existingProduct.UserID != vendorID {
		return nil, fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

	// Update only provided fields
	if req.Name != nil && *req.Name != "" {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil && *req.Description != "" {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil && *req.Price > 0 {
		existingProduct.Price = *req.Price
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	updatedProduct, err := ps.repo.UpdateProduct(ctx, existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return mapProductToResponse(updatedProduct), nil
}

// DeleteProduct deletes a product
func (ps *ProductService) DeleteProduct(ctx context.Context, productID string, vendorID string) error {
	if productID == "" || vendorID == "" {
		return fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	// Verify product exists and belongs to vendor
	product, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

	return ps.repo.DeleteProduct(ctx, productID)
}

// GetActiveProducts retrieves all active products
func (ps *ProductService) GetActiveProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetActiveProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// GetActiveUserProducts retrieves all active products for a user
func (ps *ProductService) GetActiveUserProducts(ctx context.Context, userID string) ([]*dto.ProductResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetActiveProductsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products for user: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// ToggleProductStatus toggles a product's active status
func (ps *ProductService) ToggleProductStatus(ctx context.Context, productID string, vendorID string, isActive bool) (*dto.ProductResponse, error) {
	if productID == "" || vendorID == "" {
		return nil, fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	product, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return nil, fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

	product.IsActive = isActive

	updated, err := ps.repo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	return mapProductToResponse(updated), nil
}

// SearchProducts searches for products by name or description
func (ps *ProductService) SearchProducts(ctx context.Context, searchTerm string) ([]*dto.ProductResponse, error) {
	if searchTerm == "" {
		return nil, fmt.Errorf("search term cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.SearchProducts(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// GetProductsByPriceRange retrieves products within a price range
func (ps *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*dto.ProductResponse, error) {
	if minPrice < 0 || maxPrice < 0 || minPrice > maxPrice {
		return nil, fmt.Errorf("invalid price range")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetProductsByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// Helper function to map product model to response DTO
func mapProductToResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		UserID:      product.UserID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

// Helper function to map multiple products to response DTOs
func mapProductsToResponse(products []*models.Product) []*dto.ProductResponse {
	if len(products) == 0 {
		return []*dto.ProductResponse{}
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = mapProductToResponse(product)
	}
	return responses
}

// GetProductsByUserID retrieves all products for a user (including inactive ones)
func (ps *ProductService) GetProductsByUserID(ctx context.Context, userID string) ([]*dto.ProductResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetProductsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return mapProductsToResponse(products), nil
}

// GetActiveProductsByUserID retrieves only active products for a user
func (ps *ProductService) GetActiveProductsByUserID(ctx context.Context, userID string) ([]*dto.ProductResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	products, err := ps.repo.GetProductsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Filter for active products only
	var activeProducts []*models.Product
	for _, product := range products {
		if product.IsActive {
			activeProducts = append(activeProducts, product)
		}
	}

	return mapProductsToResponse(activeProducts), nil
}
