package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/models"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	ApproveVendor(id string) error
	GetByStoreSlug(slug string) (*models.User, error)
	UpdateStoreSettings(userID, storeName, storeSlug, bio, whatsapp string) error
	GetApprovedVendors() ([]models.User, error)
}

type AuthService struct {
	userRepo  UserRepository
	jwtSecret string
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) SignUp(req dto.SignUpRequest) (*dto.AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// generate slug and ensure uniqueness
	baseSlug := utils.GenerateSlug(req.StoreName)
	slug := baseSlug
	i := 1
	for {
		if existing, _ := s.userRepo.GetByStoreSlug(slug); existing == nil {
			break
		}
		i++
		slug = baseSlug + "-" + strconv.Itoa(i)
	}

	user := &models.User{
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   hash,
		WhatsappNumber: req.WhatsappNumber,
		Username:       req.Username,
		Bio:            req.Bio,
		StoreName:      req.StoreName,
		StoreSlug:      slug,
		Role:           "vendor",
		IsActive:       false,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	authUser := dto.AuthUser{
		ID:             createdUser.ID,
		Name:           createdUser.Name,
		Email:          createdUser.Email,
		Username:       createdUser.Username,
		Role:           createdUser.Role,
		StoreName:      createdUser.StoreName,
		StoreSlug:      createdUser.StoreSlug,
		WhatsappNumber: createdUser.WhatsappNumber,
		Bio:            createdUser.Bio,
	}

	return &dto.AuthResponse{
		Token: "", User: authUser,
	}, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, utils.ErrAccountNotActive
	}

	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	token, err := utils.GenerateJwt(user)
	if err != nil {
		return nil, err
	}

	authUser := dto.AuthUser{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		StoreName:      user.StoreName,
		StoreSlug:      user.StoreSlug,
		Role:           user.Role,
		WhatsappNumber: user.WhatsappNumber,
		Bio:            user.Bio,
	}

	return &dto.AuthResponse{
		Token: token,
		User:  authUser,
	}, nil
}

func (s *AuthService) GetMyProfile(id string) (*dto.AuthUser, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return nil, utils.ErrUnauthorized
		}
		return nil, err
	}

	authUser := &dto.AuthUser{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		StoreName:      user.StoreName,
		StoreSlug:      user.StoreSlug,
		Role:           user.Role,
		WhatsappNumber: user.WhatsappNumber,
		Bio:            user.Bio,
	}

	return authUser, nil
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) GetVendorBySlug(slug string) (*models.User, error) {
	return s.userRepo.GetByStoreSlug(slug)
}

func (s *AuthService) UpdateVendorStore(ctx context.Context, userID string, req dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	storeName := user.StoreName
	if req.StoreName != nil {
		storeName = *req.StoreName
	}

	username := user.Username
	if req.Username != nil {
		username = *req.Username
	}

	bio := user.Bio
	if req.Bio != nil {
		bio = *req.Bio
	}

	whatsapp := user.WhatsappNumber
	if req.WhatsappNumber != nil {
		whatsapp = *req.WhatsappNumber
	}

	email := user.Email
	if req.Email != nil {
		email = *req.Email
	}

	storeSlug := utils.GenerateSlug(storeName)
	if storeSlug == "" {
		storeSlug = user.StoreSlug
	}

	err = s.userRepo.UpdateStoreSettings(userID, storeName, storeSlug, bio, whatsapp)
	if err != nil {
		return nil, err
	}

	return &dto.StoreResponse{
		ID:             user.ID,
		Name:           storeName,
		Slug:           storeSlug,
		Username:       username,
		Bio:            bio,
		WhatsappNumber: whatsapp,
		Email:          email,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) GetAllActiveVendors(page, pageSize int) ([]*models.User, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	vendors, err := s.userRepo.GetApprovedVendors()
	if err != nil {
		return nil, err
	}

	// convert to []*models.User
	vendorPtrs := make([]*models.User, len(vendors))
	for i := range vendors {
		vendorPtrs[i] = &vendors[i]
	}

	// pagination
	start := (page - 1) * pageSize
	if start >= len(vendorPtrs) {
		return []*models.User{}, nil
	}
	end := start + pageSize
	if end > len(vendorPtrs) {
		end = len(vendorPtrs)
	}

	return vendorPtrs[start:end], nil
}

func (s *AuthService) SearchVendors(searchTerm string) ([]*models.User, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return []*models.User{}, nil
	}

	vendors, err := s.userRepo.GetApprovedVendors()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(searchTerm)
	var results []*models.User
	for i := range vendors {
		v := &vendors[i]
		if strings.Contains(strings.ToLower(v.StoreName), lower) ||
			strings.Contains(strings.ToLower(v.Username), lower) ||
			strings.Contains(strings.ToLower(v.Name), lower) {
			results = append(results, v)
		}
	}

	return results, nil
}
