package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	_ "github.com/falasefemi2/vendorhub/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/falasefemi2/vendorhub/internal/config"
	"github.com/falasefemi2/vendorhub/internal/db"
	"github.com/falasefemi2/vendorhub/internal/handlers"
	"github.com/falasefemi2/vendorhub/internal/middleware"
	"github.com/falasefemi2/vendorhub/internal/repository"
	"github.com/falasefemi2/vendorhub/internal/service"
)

// @title VendorHub API
// @version 1.0
// @description This is a sample server for a vendor hub.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
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
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

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
		})
	})

	// Vendor public routes
	r.Route("/vendors", func(r chi.Router) {
		r.Get("/{id}/products", productHandler.GetVendorProducts)
		r.Get("/{id}/products/active", productHandler.GetActiveProducts)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
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
