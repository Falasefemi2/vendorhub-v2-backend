package service

import (
	"errors"

	"github.com/falasefemi2/vendorhub/internal/models"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type AdminRepository interface {
	GetByID(id string) (*models.User, error)
	ApproveVendor(id string) error
	GetPendingVendors() ([]models.User, error)
	GetApprovedVendors() ([]models.User, error)
}

type AdminService struct {
	userRepo AdminRepository
}

func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{userRepo: repo}
}

func (s *AdminService) ApproveVendor(adminID, vendorID string) error {
	admin, err := s.userRepo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return utils.ErrUnauthorized
		}
		return err
	}
	if admin.Role != "admin" {
		return utils.ErrUnauthorized
	}
	vendor, err := s.userRepo.GetByID(vendorID)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return utils.ErrInvalidOperation
		}
		return err
	}
	if vendor.Role != "vendor" {
		return utils.ErrInvalidOperation
	}
	if vendor.IsActive {
		return utils.ErrInvalidOperation
	}
	if err := s.userRepo.ApproveVendor(vendorID); err != nil {
		return err
	}
	return nil
}

func (s *AdminService) ListPendingVendors(adminID string) ([]models.User, error) {
	admin, err := s.userRepo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return nil, utils.ErrUnauthorized
		}
		return nil, err
	}
	if admin.Role != "admin" {
		return nil, utils.ErrUnauthorized
	}
	return s.userRepo.GetPendingVendors()
}

func (s *AdminService) ListApprovedVendors(adminID string) ([]models.User, error) {
	admin, err := s.userRepo.GetByID(adminID)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return nil, utils.ErrUnauthorized
		}
		return nil, err
	}
	if admin.Role != "admin" {
		return nil, utils.ErrUnauthorized
	}
	return s.userRepo.GetApprovedVendors()
}
