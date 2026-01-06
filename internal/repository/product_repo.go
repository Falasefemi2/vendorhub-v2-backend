package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/vendorhub/internal/models"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	product.ID = uuid.New().String()

	query := `
	INSERT INTO products (
		id, user_id, name, description, price, is_active
	) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, user_id, name, description, price, is_active, created_at, updated_at
	`

	err := pr.pool.QueryRow(
		ctx,
		query,
		product.ID,
		product.UserID,
		product.Name,
		product.Description,
		product.Price,
		product.IsActive,
	).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, productID string) (*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE id = $1
	`

	product := &models.Product{}

	err := pr.pool.QueryRow(ctx, query, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	UPDATE products
	SET name = $2, description = $3, price = $4, is_active = $5, updated_at = NOW()
	WHERE id = $1
	RETURNING id, user_id, name, description, price, is_active, created_at, updated_at
	`

	err := pr.pool.QueryRow(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.IsActive,
	).Scan(
		&product.ID,
		&product.UserID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `DELETE FROM products WHERE id = $1`

	result, err := pr.pool.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (pr *ProductRepository) GetProductsByUserID(ctx context.Context, userID string) ([]*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := pr.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

func (pr *ProductRepository) GetActiveProductsByUserID(ctx context.Context, userID string) ([]*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE user_id = $1 AND is_active = true
	ORDER BY created_at DESC
	`

	rows, err := pr.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products for user: %w", err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

func (pr *ProductRepository) GetActiveProducts(ctx context.Context) ([]*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE is_active = true
	ORDER BY created_at DESC
	`

	rows, err := pr.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

func (pr *ProductRepository) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE is_active = true AND price BETWEEN $1 AND $2
	ORDER BY price ASC
	`

	rows, err := pr.pool.Query(ctx, query, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by price range: %w", err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

func (pr *ProductRepository) SearchProducts(ctx context.Context, searchTerm string) ([]*models.Product, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, user_id, name, description, price, is_active, created_at, updated_at
	FROM products
	WHERE is_active = true AND (name ILIKE $1 OR description ILIKE $1)
	ORDER BY created_at DESC
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := pr.pool.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// CreateProductImage inserts a new product image into the database
func (pr *ProductRepository) CreateProductImage(ctx context.Context, image *models.ProductImage) (*models.ProductImage, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	image.ID = uuid.New().String()

	query := `
	INSERT INTO product_images (id, product_id, image_url, position)
	VALUES ($1, $2, $3, $4)
	RETURNING id, product_id, image_url, position, created_at
	`

	err := pr.pool.QueryRow(
		ctx,
		query,
		image.ID,
		image.ProductID,
		image.ImageURL,
		image.Position,
	).Scan(
		&image.ID,
		&image.ProductID,
		&image.ImageURL,
		&image.Position,
		&image.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product image: %w", err)
	}

	return image, nil
}

// GetProductImages retrieves all images for a product
func (pr *ProductRepository) GetProductImages(ctx context.Context, productID string) ([]*models.ProductImage, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, product_id, image_url, position, created_at
	FROM product_images
	WHERE product_id = $1
	ORDER BY position ASC
	`

	rows, err := pr.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	defer rows.Close()

	var images []*models.ProductImage

	for rows.Next() {
		image := &models.ProductImage{}
		err := rows.Scan(
			&image.ID,
			&image.ProductID,
			&image.ImageURL,
			&image.Position,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product image: %w", err)
		}
		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product images: %w", err)
	}

	return images, nil
}

// DeleteProductImage removes a product image from the database
func (pr *ProductRepository) DeleteProductImage(ctx context.Context, imageID string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `DELETE FROM product_images WHERE id = $1`

	result, err := pr.pool.Exec(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete product image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product image not found")
	}

	return nil
}

// UpdateProductImagePosition updates the position of a product image
func (pr *ProductRepository) UpdateProductImagePosition(ctx context.Context, imageID string, position int) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `UPDATE product_images SET position = $1 WHERE id = $2`

	result, err := pr.pool.Exec(ctx, query, position, imageID)
	if err != nil {
		return fmt.Errorf("failed to update image position: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product image not found")
	}

	return nil
}

// GetProductImage retrieves a single product image by ID
func (pr *ProductRepository) GetProductImage(ctx context.Context, imageID string) (*models.ProductImage, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	query := `
	SELECT id, product_id, image_url, position, created_at
	FROM product_images
	WHERE id = $1
	`

	image := &models.ProductImage{}

	err := pr.pool.QueryRow(ctx, query, imageID).Scan(
		&image.ID,
		&image.ProductID,
		&image.ImageURL,
		&image.Position,
		&image.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product image not found")
		}
		return nil, fmt.Errorf("failed to get product image: %w", err)
	}

	return image, nil
}
