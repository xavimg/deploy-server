package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
)

// UserService is a contract about something that this service can do
type UserService interface {
	Profile(userID interface{}) (entity.User, error)
	Update(user dto.UserUpdateDTO, userID interface{}, newInfo dto.UserUpdateDTO) (entity.User, error)
	DeleteAccount(userID float64) error
	VerifyAccount(email string) entity.User
	CheckRole(id interface{}) entity.TypeUser
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Profile(userID interface{}) (entity.User, error) {
	user, err := service.userRepository.ProfileUser(userID)
	if err != nil {
		return user, err
	}

	return user, err
}

func (service *userService) DeleteAccount(userID float64) error {
	if err := service.userRepository.DeleteAccount(userID); err != nil {
		return err
	}
	return nil
}

func (service *userService) Update(dataUser dto.UserUpdateDTO, userID interface{}, newInfo dto.UserUpdateDTO) (entity.User, error) {
	passToUpdate := entity.User{}

	err := smapping.FillStruct(&passToUpdate, smapping.MapFields(&dataUser))

	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}

	res, err := service.userRepository.UpdateUser(passToUpdate, userID, newInfo)
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}

func (service *userService) VerifyAccount(email string) entity.User {

	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	return user
}

func (service *userService) CheckRole(id interface{}) entity.TypeUser {

	typeUser, err := service.userRepository.CheckRole(id)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	return typeUser
}
