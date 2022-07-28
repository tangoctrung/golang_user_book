package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/tangoctrung/golang_api_v2/dto"
	"github.com/tangoctrung/golang_api_v2/entity"
	"github.com/tangoctrung/golang_api_v2/repository"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (c *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updateUser := c.userRepository.UpdateUser(userToUpdate)
	return updateUser
}

func (c *userService) Profile(userID string) entity.User {
	return c.userRepository.ProfileUser(userID)
}
