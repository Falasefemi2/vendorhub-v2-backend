package dto

type StoreResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	WhatsappNumber string `json:"whatsapp_number"`
	Email          string `json:"email"`
	CreatedAt      string `json:"created_at"`
}

type StoreDetailsResponse struct {
	Store    *StoreResponse     `json:"store"`
	Products []*ProductResponse `json:"products"`
	StoreURL string             `json:"store_url"`
}

type UpdateStoreRequest struct {
	StoreName      *string `json:"store_name"`
	Username       *string `json:"username"`
	Bio            *string `json:"bio"`
	WhatsappNumber *string `json:"whatsapp_number"`
	Email          *string `json:"email"`
}
