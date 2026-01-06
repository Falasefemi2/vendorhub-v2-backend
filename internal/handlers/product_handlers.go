package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/models"
	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/storage"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type ProductHandler struct {
	service *service.ProductService
	storage storage.Storage
}

func NewProductHandler(service *service.ProductService, storage storage.Storage) *ProductHandler {
	return &ProductHandler{service: service, storage: storage}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product for the authenticated vendor
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body body dto.CreateProductRequest true "Create Product Request"
// @Success      201  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products [post]
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

// GetProduct godoc
// @Summary      Get a product by ID
// @Description  Retrieves a single product by its ID with images
// @Tags         Products
// @Produce      json
// @Param        id   query      string  true  "Product ID"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products [get]
func (ph *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "product id is required")
		return
	}

	response, err := ph.service.GetProductWithImages(r.Context(), productID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetUserProducts godoc
// @Summary      Get authenticated vendor's products
// @Description  Retrieves all products for the currently authenticated vendor
// @Tags         Products
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   dto.ProductResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/my [get]
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

// GetVendorProducts godoc
// @Summary      Get a vendor's products
// @Description  Retrieves all products for a specific vendor
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Vendor ID"
// @Success      200  {array}   dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /vendors/{id}/products [get]
func (ph *ProductHandler) GetVendorProducts(w http.ResponseWriter, r *http.Request) {
	vendorID := chi.URLParam(r, "id")
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

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product for the authenticated vendor
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Product ID"
// @Param        body body      dto.UpdateProductRequest true "Update Product Request"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/{id} [put]
func (ph *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
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

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes a product for the authenticated vendor
// @Tags         Products
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/{id} [delete]
func (ph *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
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

// GetActiveProducts godoc
// @Summary      Get active products
// @Description  Retrieves all active products
// @Tags         Products
// @Produce      json
// @Success      200  {array}   dto.ProductResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/active [get]
func (ph *ProductHandler) GetActiveProducts(w http.ResponseWriter, r *http.Request) {
	responses, err := ph.service.GetActiveProducts(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// GetActiveUserProducts godoc
// @Summary      Get active products for a vendor
// @Description  Retrieves all active products for a specific vendor
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Vendor ID"
// @Success      200  {array}   dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /vendors/{id}/products/active [get]
func (ph *ProductHandler) GetActiveUserProducts(w http.ResponseWriter, r *http.Request) {
	vendorID := chi.URLParam(r, "id")
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

// ToggleProductStatus godoc
// @Summary      Toggle a product's active status
// @Description  Toggles a product's active status for the authenticated vendor
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Product ID"
// @Param        body body      dto.ToggleProductStatusRequest  true  "Toggle Status Request"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/{id}/status [put]
func (ph *ProductHandler) ToggleProductStatus(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
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
	var req dto.ToggleProductStatusRequest
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

// SearchProducts godoc
// @Summary      Search for products
// @Description  Searches for products by a search term
// @Tags         Products
// @Produce      json
// @Param        q    query     string  true  "Search Term"
// @Success      200  {array}   dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/search [get]
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

// GetProductsByPriceRange godoc
// @Summary      Get products by price range
// @Description  Retrieves products within a specified price range
// @Tags         Products
// @Produce      json
// @Param        min  query     number  true  "Minimum Price"
// @Param        max  query     number  true  "Maximum Price"
// @Success      200  {array}   dto.ProductResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/price [get]
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

// UploadProductImage godoc
// @Summary      Upload an image for a product
// @Description  Uploads a new image for a product
// @Tags         ProductImages
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        productId path string true "Product ID"
// @Param        image formData file true "Image file"
// @Param        position formData integer false "Image position"
// @Success      201  {object}  dto.UploadProductImageResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /products/{productId}/images [post]
func (ph *ProductHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
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
		utils.WriteError(w, http.StatusForbidden, "only vendors can upload product images")
		return
	}

	// Parse multipart form with max 10MB size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	// Get image file
	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	// Get position from form (optional, defaults to 0)
	position := 0
	if posStr := r.FormValue("position"); posStr != "" {
		if pos, err := strconv.Atoi(posStr); err == nil && pos >= 0 {
			position = pos
		}
	}

	// Save file
	filename, err := ph.storage.SaveFile(r.Context(), handler)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create image record
	imageReq := &dto.UploadProductImageRequest{Position: position}
	productImage := &models.ProductImage{
		ImageURL: filename,
	}

	response, err := ph.service.CreateProductImage(r.Context(), productID, vendorID, imageReq, productImage)
	if err != nil {
		// Clean up file if database operation fails
		ph.storage.DeleteFile(r.Context(), filename)
		utils.HandleServiceError(w, err)
		return
	}

	// Add full URL to response
	response.ImageURL = ph.storage.GetURL(filename)

	utils.WriteJSON(w, http.StatusCreated, response)
}

// DeleteProductImage godoc
// @Summary      Delete a product image
// @Description  Deletes an image from a product
// @Tags         ProductImages
// @Produce      json
// @Security     ApiKeyAuth
// @Param        imageId path string true "Image ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /images/{imageId} [delete]
func (ph *ProductHandler) DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "imageId")
	if imageID == "" {
		utils.WriteError(w, http.StatusBadRequest, "image id is required")
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
		utils.WriteError(w, http.StatusForbidden, "only vendors can delete product images")
		return
	}

	err = ph.service.DeleteProductImage(r.Context(), imageID, vendorID)
	if err != nil {
		if err.Error() == "unauthorized: image does not belong to this vendor" {
			utils.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "image deleted successfully"})
}

// UpdateProductImagePosition godoc
// @Summary      Update product image position
// @Description  Changes the position/order of a product image
// @Tags         ProductImages
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        imageId path string true "Image ID"
// @Param        body body dto.UploadProductImageRequest true "Position Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /images/{imageId}/position [put]
func (ph *ProductHandler) UpdateProductImagePosition(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "imageId")
	if imageID == "" {
		utils.WriteError(w, http.StatusBadRequest, "image id is required")
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
		utils.WriteError(w, http.StatusForbidden, "only vendors can update product images")
		return
	}

	var req dto.UploadProductImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	err = ph.service.UpdateProductImagePosition(r.Context(), imageID, vendorID, req.Position)
	if err != nil {
		if err.Error() == "unauthorized: image does not belong to this vendor" {
			utils.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "image position updated successfully"})
}
