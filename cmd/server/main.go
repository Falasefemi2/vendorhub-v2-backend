package main

//go:generate swag init --parseDependency --parseInternal -g cmd/server/main.go -d ./,./internal/handlers
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger"

	docs "github.com/falasefemi2/vendorhub/docs"

	"github.com/falasefemi2/vendorhub/internal/config"
	"github.com/falasefemi2/vendorhub/internal/db"
	"github.com/falasefemi2/vendorhub/internal/handlers"
	"github.com/falasefemi2/vendorhub/internal/middleware"
	"github.com/falasefemi2/vendorhub/internal/repository"
	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/storage"
)

// @title VendorHub API
// @version 1.0
// @description This is a sample server for a vendor hub.
// @host vendorhub-v2-backend-2.onrender.com
// @schemes https
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.Load()
	connString := config.GetDBURL()
	ctx := context.Background()

	pool, err := db.ConnectAndMigrate(ctx, connString)
	if err != nil {
		panic(fmt.Errorf("failed to migrate: %w", err))
	}
	defer pool.Close()

	fmt.Println("Database ready")

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	authHandler := handlers.NewAuthHandler(authService)

	adminService := service.NewAdminService(userRepo)
	adminHandler := handlers.NewAdminHandler(adminService)

	productRepo := repository.NewProductRepository(pool)

	// Initialize Supabase storage
	supabaseURL := config.GetSupabaseURL()
	supabaseKey := config.GetSupabaseKey()
	supabaseBucket := config.GetSupabaseBucket()

	supabaseStorage, err := storage.NewSupabaseStorage(supabaseURL, supabaseKey, supabaseBucket)
	if err != nil {
		fmt.Printf("error: failed to initialize Supabase storage: %v\n", err)
		panic(fmt.Errorf("failed to initialize Supabase storage: %w", err))
	}

	productService := service.NewProductService(productRepo, supabaseStorage)
	productHandler := handlers.NewProductHandler(productService, supabaseStorage)

	storeHandler := handlers.NewStoreHandler(authService, productService)

	// Configure Swagger host/schemes at runtime so local testing uses localhost:8080
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		// running locally
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		// remove scheme and trailing slash
		host := strings.TrimPrefix(baseURL, "https://")
		host = strings.TrimPrefix(host, "http://")
		host = strings.TrimSuffix(host, "/")
		docs.SwaggerInfo.Host = host
		if strings.HasPrefix(baseURL, "https://") {
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else {
			docs.SwaggerInfo.Schemes = []string{"http"}
		}
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing health check response: %v", err)
		}
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Use(middleware.AdminOnly)

		r.Post("/vendors/{id}/approve", adminHandler.ApproveVendor)
		r.Get("/vendors/pending", adminHandler.ListPendingVendors)
		r.Get("/vendors/approved", adminHandler.ListApprovedVendors)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/me", authHandler.GetMyProfile)
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/active", productHandler.GetActiveProducts)
		r.Get("/search", productHandler.SearchProducts)
		r.Get("/price", productHandler.GetProductsByPriceRange)
		r.Get("/", productHandler.GetProduct)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)

			// Vendor-only operations
			r.Post("/", productHandler.CreateProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
			r.Put("/{id}/status", productHandler.ToggleProductStatus)
			r.Get("/my", productHandler.GetUserProducts)

			// Product image operations
			r.Post("/{productId}/images", productHandler.UploadProductImage)
		})
	})

	// Image management routes (vendor-only)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Route("/images", func(r chi.Router) {
			r.Delete("/{imageId}", productHandler.DeleteProductImage)
			r.Put("/{imageId}/position", productHandler.UpdateProductImagePosition)
		})
	})

	r.Route("/stores", func(r chi.Router) {
		// Public store endpoints
		// GET /stores - All vendors with stores
		r.Get("/", storeHandler.GetAllStores)

		// GET /stores/search?q=pizza - Search vendors
		r.Get("/search", storeHandler.SearchStores)

		// GET /stores/vendor?id={vendorId} - Get vendor's store by ID
		r.Get("/vendor", storeHandler.GetStoreByVendorID)

		// WHATSAPP SHAREABLE LINK
		// GET /stores/@{store-slug} - Get vendor store + products by slug
		// Example: GET /stores/@pizzahut-lagos
		r.Get("/{slug}", storeHandler.GetStoreBySlug)

		// Protected store endpoints (vendor only)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)

			// GET /stores/my - Get authenticated vendor's store with products
			r.Get("/my", storeHandler.GetMyStore)

			// PUT /stores/my - Update vendor's store info
			r.Put("/my", storeHandler.UpdateMyStore)
		})
	})

	// Vendor public routes
	r.Route("/vendors", func(r chi.Router) {
		r.Get("/{id}/products", productHandler.GetVendorProducts)
		r.Get("/{id}/products/active", productHandler.GetActiveProducts)
	})

	// Build CORS allowed origins
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"https://vendorhub-v2-frontend.vercel.app",
	}

	if prodOrigins := os.Getenv("ALLOWED_ORIGINS"); prodOrigins != "" {
		origins := strings.Split(prodOrigins, ",")
		for _, origin := range origins {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      c.Handler(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
