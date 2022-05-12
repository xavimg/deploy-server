package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
)

type AdminService interface {
	CreateAdmin(admin dto.RegisterDTO) (entity.User, error)
	ListAllUsers() ([]entity.User, error)
	ListAllUsersByActive() ([]entity.User, error)
	ListAllUsersByTypeAdmin() ([]entity.User, error)
	ListAllUsersByTypeUser() ([]entity.User, error)
	BanUser(userID string) error
	UnbanUser(userID string) error
	NewFeature(feature dto.FeatureDTO) (entity.Feature, error)
}

type adminService struct {
	adminRepository repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepository: adminRepo,
	}
}

func (service *adminService) CreateAdmin(user dto.RegisterDTO) (entity.User, error) {
	adminToCreate := entity.User{}

	if err := smapping.FillStruct(&adminToCreate, smapping.MapFields(&user)); err != nil {
		log.Fatalf("Failed map %v", err)
		return entity.User{}, err
	}

	res, err := service.adminRepository.InsertAdmin(adminToCreate)
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}

func (service *adminService) ListAllUsers() ([]entity.User, error) {
	users, err := service.adminRepository.ListAllUsers()
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (service *adminService) ListAllUsersByActive() ([]entity.User, error) {
	users, err := service.adminRepository.ListAllUsersByActive()
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (service *adminService) ListAllUsersByTypeAdmin() ([]entity.User, error) {
	users, err := service.adminRepository.ListAllUsersByTypeAdmin()
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (service *adminService) ListAllUsersByTypeUser() ([]entity.User, error) {
	users, err := service.adminRepository.ListAllUsersByTypeUser()
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (service *adminService) BanUser(userID string) error {
	if err := service.adminRepository.BanUser(userID); err != nil {
		return err
	}

	return nil
}

func (service *adminService) UnbanUser(userID string) error {
	if err := service.adminRepository.UnbanUser(userID); err != nil {
		return err
	}

	return nil
}

func (service *adminService) NewFeature(feature dto.FeatureDTO) (entity.Feature, error) {
	featureToCreate := entity.Feature{}

	if err := smapping.FillStruct(&featureToCreate, smapping.MapFields(&feature)); err != nil {
		return entity.Feature{}, err
	}

	res, err := service.adminRepository.NewFeature(featureToCreate)
	if err != nil {
		return entity.Feature{}, err
	}

	return res, nil
}
