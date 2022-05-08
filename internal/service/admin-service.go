package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
)

type AdminService interface {
	CreateAdmin(admin dto.RegisterDTO) entity.User
	ListAllUsers() []entity.User
	ListAllUsersByActive() []entity.User
	ListAllUsersByTypeAdmin() []entity.User
	ListAllUsersByTypeUser() []entity.User
	BanUser(userID string)
	UnbanUser(userID string)
	NewFeature(feature dto.FeatureDTO) entity.Feature
}

type adminService struct {
	adminRepository repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepository: adminRepo,
	}
}

func (service *adminService) CreateAdmin(user dto.RegisterDTO) entity.User {
	adminToCreate := entity.User{}

	err := smapping.FillStruct(&adminToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	res := service.adminRepository.InsertAdmin(adminToCreate)

	return res
}

func (service *adminService) ListAllUsers() []entity.User {
	users := service.adminRepository.ListAllUsers()
	return users
}

func (service *adminService) ListAllUsersByActive() []entity.User {
	users := service.adminRepository.ListAllUsersByActive()
	return users
}

func (service *adminService) ListAllUsersByTypeAdmin() []entity.User {
	users := service.adminRepository.ListAllUsersByTypeAdmin()
	return users
}

func (service *adminService) ListAllUsersByTypeUser() []entity.User {
	users := service.adminRepository.ListAllUsersByTypeUser()
	return users
}

func (service *adminService) BanUser(userID string) {
	service.adminRepository.BanUser(userID)
}

func (service *adminService) UnbanUser(userID string) {
	service.adminRepository.UnbanUser(userID)
}

func (service *adminService) NewFeature(feature dto.FeatureDTO) entity.Feature {
	featureToCreate := entity.Feature{}

	err := smapping.FillStruct(&featureToCreate, smapping.MapFields(&feature))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	res := service.adminRepository.NewFeature(featureToCreate)

	return res
}
