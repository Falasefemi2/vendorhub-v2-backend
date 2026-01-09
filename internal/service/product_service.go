package service

import (
	"context"
	"fmt"
	"time"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/models"
	"github.com/falasefemi2/vendorhub/internal/repository"
	"github.com/falasefemi2/vendorhub/internal/storage"
)

type ProductService struct {
	repo    *repository.ProductRepository
	storage storage.Storage
}

func NewProductService(repo *repository.ProductRepository, storage storage.Storage) *ProductService {
	return &ProductService{repo: repo, storage: storage}
}

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

	responses := mapProductsToResponse(products)
	// Enrich with images
	ps.enrichProductResponsesWithImages(ctx, responses)
	return responses, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, productID string, vendorID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if productID == "" || vendorID == "" {
		return nil, fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	existingProduct, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if existingProduct.UserID != vendorID {
		return nil, fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

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

func (ps *ProductService) DeleteProduct(ctx context.Context, productID string, vendorID string) error {
	if productID == "" || vendorID == "" {
		return fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	product, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

	return ps.repo.DeleteProduct(ctx, productID)
}

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

	responses := mapProductsToResponse(products)
	// Enrich with images
	ps.enrichProductResponsesWithImages(ctx, responses)
	return responses, nil
}

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

	responses := mapProductsToResponse(products)
	// Enrich with images
	ps.enrichProductResponsesWithImages(ctx, responses)
	return responses, nil
}

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

	responses := mapProductsToResponse(products)
	// Enrich with images
	ps.enrichProductResponsesWithImages(ctx, responses)
	return responses, nil
}

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

	responses := mapProductsToResponse(products)
	// Enrich with images
	ps.enrichProductResponsesWithImages(ctx, responses)
	return responses, nil
}

func mapProductToResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID,
		UserID:      product.UserID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		IsActive:    product.IsActive,
		Images:      []*dto.ProductImageResponse{},
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

// enrichProductResponseWithImages adds images to a product response
func (ps *ProductService) enrichProductResponseWithImages(ctx context.Context, response *dto.ProductResponse) error {
	images, err := ps.repo.GetProductImages(ctx, response.ID)
	if err != nil {
		// Don't fail if we can't get images, just return without images
		return nil
	}
	response.Images = ps.mapProductImagesToResponse(images)
	return nil
}

// enrichProductResponsesWithImages adds images to multiple product responses
func (ps *ProductService) enrichProductResponsesWithImages(ctx context.Context, responses []*dto.ProductResponse) error {
	for _, response := range responses {
		if err := ps.enrichProductResponseWithImages(ctx, response); err != nil {
			// Continue enriching even if one fails
			continue
		}
	}
	return nil
}

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

	var activeProducts []*models.Product
	for _, product := range products {
		if product.IsActive {
			activeProducts = append(activeProducts, product)
		}
	}

	return mapProductsToResponse(activeProducts), nil
}

// GetProductWithImages retrieves a product with its images
func (ps *ProductService) GetProductWithImages(ctx context.Context, productID string) (*dto.ProductResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("product ID cannot be empty")
	}

	product, err := ps.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	images, err := ps.repo.GetProductImages(ctx, productID)
	if err != nil {
		// Don't fail if we can't get images, just return product without images
		return product, nil
	}

	product.Images = ps.mapProductImagesToResponse(images)
	return product, nil
}

// CreateProductImage saves an image file and creates image record
func (ps *ProductService) CreateProductImage(ctx context.Context, productID string, vendorID string, req *dto.UploadProductImageRequest, file *models.ProductImage) (*dto.ProductImageResponse, error) {
	if productID == "" || vendorID == "" {
		return nil, fmt.Errorf("product ID and vendor ID cannot be empty")
	}

	// Verify product belongs to vendor
	product, err := ps.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return nil, fmt.Errorf("unauthorized: product does not belong to this vendor")
	}

	// Create the image record in database
	image := &models.ProductImage{
		ProductID: productID,
		ImageURL:  file.ImageURL,
		Position:  req.Position,
	}

	createdImage, err := ps.repo.CreateProductImage(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("failed to create product image: %w", err)
	}

	return ps.mapProductImageToResponse(createdImage), nil
}

// DeleteProductImage removes an image file and database record
func (ps *ProductService) DeleteProductImage(ctx context.Context, imageID string, vendorID string) error {
	if imageID == "" || vendorID == "" {
		return fmt.Errorf("image ID and vendor ID cannot be empty")
	}

	// Get the image to find the product
	image, err := ps.repo.GetProductImage(ctx, imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Verify product belongs to vendor
	product, err := ps.repo.GetProductByID(ctx, image.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return fmt.Errorf("unauthorized: image does not belong to this vendor")
	}

	// Delete file from storage
	if err := ps.storage.DeleteFile(ctx, image.ImageURL); err != nil {
		// Log error but don't fail the whole operation
		fmt.Printf("warning: failed to delete image file %s: %v\n", image.ImageURL, err)
	}

	// Delete image record from database
	return ps.repo.DeleteProductImage(ctx, imageID)
}

// UpdateProductImagePosition changes the position of an image
func (ps *ProductService) UpdateProductImagePosition(ctx context.Context, imageID string, vendorID string, newPosition int) error {
	if imageID == "" || vendorID == "" {
		return fmt.Errorf("image ID and vendor ID cannot be empty")
	}

	if newPosition < 0 {
		return fmt.Errorf("image position cannot be negative")
	}

	// Get the image to find the product
	image, err := ps.repo.GetProductImage(ctx, imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Verify product belongs to vendor
	product, err := ps.repo.GetProductByID(ctx, image.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.UserID != vendorID {
		return fmt.Errorf("unauthorized: image does not belong to this vendor")
	}

	return ps.repo.UpdateProductImagePosition(ctx, imageID, newPosition)
}

// mapProductImageToResponse maps a models.ProductImage to a DTO, converting stored
// filename to a full accessible URL using the configured storage.
func (ps *ProductService) mapProductImageToResponse(image *models.ProductImage) *dto.ProductImageResponse {
	return &dto.ProductImageResponse{
		ID: image.ID,
		// ImageURL: ps.storage.GetURL(image.ImageURL),
		ImageURL: image.ImageURL,
		Position: image.Position,
	}
}

func (ps *ProductService) mapProductImagesToResponse(images []*models.ProductImage) []*dto.ProductImageResponse {
	if len(images) == 0 {
		return []*dto.ProductImageResponse{}
	}

	responses := make([]*dto.ProductImageResponse, len(images))
	for i, image := range images {
		responses[i] = ps.mapProductImageToResponse(image)
	}
	return responses
}
