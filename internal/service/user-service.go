package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/mashingan/smapping"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
)

// UserService is a contract about something that this service can do
type UserService interface {
	Profile(userID interface{}) (*entity.User, error)
	Update(user dto.UserUpdateDTO, userID interface{}, newInfo dto.UserUpdateDTO) (entity.User, error)
	DeleteAccount(userID float64) error
	VerifyAccount(email string) entity.User
	CheckRole(id interface{}) entity.TypeUser

	AddFriend(id interface{}, friend *entity.User) error
	ShowFriendlist(id interface{}) ([]dto.Friend, error)
	RemoveFriend(idUser interface{}, id uint64) error
	IsFriend(id interface{}) (bool, error)
	SendMessage(message dto.MessageDTO) error
	ListMessages(id interface{}) ([]*dto.NotificationMessageDTO, error)
	MessageDetail(idMessage int, idAuth interface{}) (*dto.BodyMessageDTO, error)
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

func (service *userService) Profile(userID interface{}) (*entity.User, error) {
	user, err := service.userRepository.ProfileUser(userID)
	if err != nil {
		return nil, err
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

func (s *userService) AddFriend(id interface{}, friend *entity.User) error {
	friendToAdd := dto.Friend{
		ID:    friend.ID,
		Name:  friend.Name,
		Email: friend.Email,
	}

	if err := s.userRepository.AddFriend(id, &friendToAdd); err != nil {
		return err
	}

	return nil
}

func (s *userService) ShowFriendlist(id interface{}) ([]dto.Friend, error) {
	rows, err := s.userRepository.ShowFriendlist(id)
	if err != nil {
		return nil, err
	}

	friend := dto.Friend{}
	friends := []dto.Friend{}
	for _, v := range rows {
		if err := json.Unmarshal(v.Friendlist, &friend); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (s *userService) RemoveFriend(idUser interface{}, id uint64) error {
	rows, err := s.userRepository.ShowFriendlist(idUser)
	if err != nil {
		return err
	}

	friend := dto.Friend{}
	friends := []dto.Friend{}
	for _, v := range rows {
		if err := json.Unmarshal(v.Friendlist, &friend); err != nil {
			return err
		}
		friends = append(friends, friend)
	}

	for _, k := range friends {
		if id == k.ID {
			if err := s.userRepository.RemoveFriend(id); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("user not found in friendlist")
}

func (s *userService) IsFriend(id interface{}) (bool, error) {
	exists, err := s.userRepository.IsFriend(id)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("requested user is not a friend")
	}

	return true, nil
}

func (s *userService) SendMessage(m dto.MessageDTO) error {
	message := entity.User_Messages{
		Receiver: m.To,
		Sender:   float64(m.From),
		Detail:   m.Message,
	}

	if err := s.userRepository.SendMessage(message); err != nil {
		return err
	}

	return nil
}

func (s *userService) ListMessages(id interface{}) ([]*dto.NotificationMessageDTO, error) {
	notifications, err := s.userRepository.ListMessages(id)
	if err != nil {
		return nil, err
	}

	notificationsDTO := []*dto.NotificationMessageDTO{}
	for _, noti := range notifications {
		notificationsDTO = append(notificationsDTO, &dto.NotificationMessageDTO{
			Remitent: float64(noti.Sender),
			Tittle:   noti.Detail,
		})
	}

	return notificationsDTO, nil
}

func (s *userService) MessageDetail(id int, idAuth interface{}) (*dto.BodyMessageDTO, error) {
	message, err := s.userRepository.MessageDetail(id)
	if err != nil {
		return nil, err
	}

	if message.Sender != idAuth.(float64) {
		return nil, errors.New("el jugador authenticat no té un missatge amb aquest id")
	}

	messageDTO := dto.BodyMessageDTO{
		Detail: message.Detail,
	}

	return &messageDTO, nil
}
