package service

import (
	"log"
	"strconv"

	"github.com/mashingan/smapping"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(user dto.RegisterDTO) entity.User
	DeleteUser(id uint64) error
	VerifyCredential(email, password string) interface{}
	VerifyUserExist(userID string) interface{}
	VerifyUserActive(email string) entity.User
	FindByEmail(email string) (entity.User, error)
	IsDuplicateEmail(email string) bool
	SaveToken(user entity.User, token string)
	DeleteToken(user entity.User, s string)
	GetToken(UserID string) entity.User
	VerifyCode(email string, code int) (bool, error)
}

type authService struct {
	userRepository repository.UserRepository
	authRepository repository.AuthRepository
}

func NewAuthService(userRepo repository.UserRepository, authRepository repository.AuthRepository) AuthService {
	return &authService{
		userRepository: userRepo,
		authRepository: authRepository,
	}
}

func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}

	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	res := service.userRepository.InsertUser(userToCreate)

	return res
}

func (service *authService) DeleteUser(id uint64) error {
	if err := service.userRepository.DeleteAccount(id); err != nil {
		return err
	}
	return nil
}

func (service *authService) VerifyCredential(email, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)

	if v, ok := res.(entity.User); ok {

		comparedPassword := comparePassword(v.Password, []byte(password))

		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}

	return false
}

func (service *authService) VerifyUserExist(id string) interface{} {

	res := service.userRepository.VerifyUserExist(id)

	if v, ok := res.(entity.User); ok {

		if strconv.FormatUint(v.ID, 10) == id {
			return res
		}
		return false
	}

	return false

}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)

	return !(res.Error == nil)
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (service *authService) FindByEmail(email string) (entity.User, error) {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) SaveToken(user entity.User, token string) {
	service.userRepository.SaveToken(user, token)
}

func (service *authService) DeleteToken(user entity.User, s string) {
	service.userRepository.DeleteToken(user, s)
}

func (service *authService) GetToken(userID string) entity.User {

	return service.userRepository.GetToken(userID)
}

func (service *authService) VerifyUserActive(email string) entity.User {

	return service.userRepository.VerifyUserActive(email)

}

func (service *authService) VerifyCode(email string, code int) (bool, error) {

	exist, err := service.authRepository.VerifyCodeByEmail(email, code)
	if err != nil {
		log.Println("Error: ", err)
		return false, err
	}

	if !exist {
		log.Println("Code not valid")
		return false, err
	}

	return true, nil
}
