package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/vendorhub/internal/models"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	query := `
		INSERT INTO users (id, name, email, password_hash, whatsapp_number, username, bio, store_name, store_slug, role, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, email, username, whatsapp_number, bio, store_name, store_slug, role, is_active, created_at
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.WhatsappNumber,
		user.Username,
		user.Bio,
		user.StoreName,
		user.StoreSlug,
		user.Role,
		user.IsActive,
		user.CreatedAt,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Username, &user.WhatsappNumber, &user.Bio, &user.StoreName, &user.StoreSlug, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, name, email, password_hash, whatsapp_number, username, bio, store_name, store_slug, role, is_active, created_at
		FROM users
		WHERE email = $1
	`

	err := r.pool.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.WhatsappNumber,
		&user.Username,
		&user.Bio,
		&user.StoreName,
		&user.StoreSlug,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, name, email, password_hash, whatsapp_number, username, bio, store_name, store_slug, role, is_active, created_at
		FROM users
		WHERE id = $1
	`

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.WhatsappNumber,
		&user.Username,
		&user.Bio,
		&user.StoreName,
		&user.StoreSlug,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByStoreSlug(slug string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, name, email, password_hash, whatsapp_number, username, bio, store_name, store_slug, role, is_active, created_at
		FROM users
		WHERE store_slug = $1
	`

	err := r.pool.QueryRow(context.Background(), query, slug).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.WhatsappNumber,
		&user.Username,
		&user.Bio,
		&user.StoreName,
		&user.StoreSlug,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("store not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateStoreSettings(userID, storeName, storeSlug, bio, whatsapp string) error {
	query := `
		UPDATE users
		SET store_name = $1, store_slug = $2, bio = $3, whatsapp_number = $4
		WHERE id = $5
	`

	result, err := r.pool.Exec(context.Background(), query, storeName, storeSlug, bio, whatsapp, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) ApproveVendor(id string) error {
	query := `
		UPDATE users
		SET is_active = true
		WHERE id = $1
	`

	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) GetPendingVendors() ([]models.User, error) {
	query := `
		SELECT id, name, email, whatsapp_number, username, bio, role, is_active, created_at, store_name, store_slug
		FROM users
		WHERE role = 'vendor' AND is_active = false
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.WhatsappNumber,
			&user.Username,
			&user.Bio,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.StoreName,
			&user.StoreSlug,
		); err != nil {
			return nil, err
		}
		vendors = append(vendors, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

func (r *UserRepository) GetApprovedVendors() ([]models.User, error) {
	query := `
		SELECT id, name, email, whatsapp_number, username, bio, role, is_active, created_at, store_name, store_slug
		FROM users
		WHERE role = 'vendor' AND is_active = true
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.WhatsappNumber,
			&user.Username,
			&user.Bio,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.StoreName,
			&user.StoreSlug,
		); err != nil {
			return nil, err
		}
		vendors = append(vendors, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}
