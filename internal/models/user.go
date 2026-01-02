package models

import "time"

type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	WhatsappNumber string    `json:"whatsapp_number"`
	Username       string    `json:"username"`
	Bio            string    `json:"bio"`
	StoreName      string    `json:"store_name"`
	StoreSlug      string    `json:"store_slug"`
	Role           string    `json:"role"` // admin | vendor
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}
