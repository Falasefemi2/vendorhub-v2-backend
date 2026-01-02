package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct creates a new product
// POST /products
func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can create products")
		return
	}

	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	response, err := ph.service.CreateProduct(r.Context(), vendorID, req)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

// GetProduct retrieves a single product by ID
// GET /products?id={productId}
func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "product id is required")
		return
	}

	response, err := ph.service.GetProduct(r.Context(), productID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetUserProducts retrieves all products for authenticated vendor
// GET /products/my
func (ph *ProductHandler) GetUserProducts(w http.ResponseWriter, r *http.Request) {
	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can view their products")
		return
	}

	response, err := ph.service.GetUserProducts(r.Context(), vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetVendorProducts retrieves all products for a specific vendor
// GET /vendors/{id}/products
func (ph *ProductHandler) GetVendorProducts(w http.ResponseWriter, r *http.Request) {
	vendorID := r.URL.Query().Get("vendor_id")
	if vendorID == "" {
		utils.WriteError(w, http.StatusBadRequest, "vendor_id is required")
		return
	}

	response, err := ph.service.GetUserProducts(r.Context(), vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// UpdateProduct updates an existing product
// PUT /products?id={productId}
func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "product id is required")
		return
	}

	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can update products")
		return
	}

	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	response, err := ph.service.UpdateProduct(r.Context(), productID, vendorID, req)
	if err != nil {
		if err.Error() == "unauthorized: product does not belong to this vendor" {
			utils.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// DeleteProduct deletes a product
// DELETE /products?id={productId}
func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "product id is required")
		return
	}

	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can delete products")
		return
	}

	err = ph.service.DeleteProduct(r.Context(), productID, vendorID)
	if err != nil {
		if err.Error() == "unauthorized: product does not belong to this vendor" {
			utils.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

// GetActiveProducts retrieves all active products
// GET /products/active
func (ph *ProductHandler) GetActiveProducts(w http.ResponseWriter, r *http.Request) {
	responses, err := ph.service.GetActiveProducts(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// GetActiveUserProducts retrieves all active products for a vendor
// GET /products/active?vendor_id={vendorId}
func (ph *ProductHandler) GetActiveUserProducts(w http.ResponseWriter, r *http.Request) {
	vendorID := r.URL.Query().Get("vendor_id")
	if vendorID == "" {
		utils.WriteError(w, http.StatusBadRequest, "vendor_id is required")
		return
	}

	responses, err := ph.service.GetActiveUserProducts(r.Context(), vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// ToggleProductStatus toggles a product's active status
// PUT /products?id={productId}&status={true|false}
func (ph *ProductHandler) ToggleProductStatus(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "product id is required")
		return
	}

	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can toggle product status")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	response, err := ph.service.ToggleProductStatus(r.Context(), productID, vendorID, req.IsActive)
	if err != nil {
		if err.Error() == "unauthorized: product does not belong to this vendor" {
			utils.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// SearchProducts searches for products
// GET /products/search?q={searchTerm}
func (ph *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		utils.WriteError(w, http.StatusBadRequest, "search term is required")
		return
	}

	responses, err := ph.service.SearchProducts(r.Context(), searchTerm)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// GetProductsByPriceRange retrieves products within a price range
// GET /products/price?min={minPrice}&max={maxPrice}
func (ph *ProductHandler) GetProductsByPriceRange(w http.ResponseWriter, r *http.Request) {
	minPriceStr := r.URL.Query().Get("min")
	maxPriceStr := r.URL.Query().Get("max")

	if minPriceStr == "" || maxPriceStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "min and max price parameters are required")
		return
	}

	var minPrice, maxPrice float64
	_, err := utils.ParseFloat64(minPriceStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid min price")
		return
	}

	maxPriceVal, err := utils.ParseFloat64(maxPriceStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid max price")
		return
	}

	minPrice, _ = utils.ParseFloat64(minPriceStr)
	maxPrice = maxPriceVal

	responses, err := ph.service.GetProductsByPriceRange(r.Context(), minPrice, maxPrice)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}
