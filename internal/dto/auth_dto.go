package dto

type SignUpRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	WhatsappNumber string `json:"whatsapp_number" binding:"required"`
	Username       string `json:"username" binding:"required"`
	StoreName      string `json:"store_name" binding:"required"`
	Bio            string `json:"bio"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthUser struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	StoreName      string `json:"store_name"`
	StoreSlug      string `json:"store_slug"`
	Role           string `json:"role"`
	Bio            string `json:"bio"`
	WhatsappNumber string `json:"whatsapp_number"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}
