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
	Profile(userID string) entity.User
	Update(user dto.UserUpdateDTO, userID string, newInfo dto.UserUpdateDTO) entity.User
	VerifyAccount(email string) entity.User
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

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}

func (service *userService) Update(dataUser dto.UserUpdateDTO, userID string, newInfo dto.UserUpdateDTO) entity.User {
	passToUpdate := entity.User{}

	err := smapping.FillStruct(&passToUpdate, smapping.MapFields(&dataUser))

	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}

	res := service.userRepository.UpdateUser(passToUpdate, userID, newInfo)

	return res
}

func (service *userService) VerifyAccount(email string) entity.User {

	user, err := service.userRepository.FindByEmail(email)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	return user
}
