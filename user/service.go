package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	UserLogin(input LoginUserInput) (User, error)
	EmailAvailability(input EmailAvailabilityInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"
	newUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}
	return newUser, nil
}

func (s *service) UserLogin(input LoginUserInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User with that email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *service) EmailAvailability(input EmailAvailabilityInput) (bool, error) {
	email := input.Email
	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, err
	}

	return false, err
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// get user by ID
	user, err := s.repository.FindByID(ID)

	if err != nil {
		return user, err
	}

	// assign avatarfilename value
	user.AvatarFileName = fileLocation

	//update user with new avatarfilename
	updatedUser, err := s.repository.UpdateUser(user)

	if err != nil {
		return updatedUser, err
	}
	return updatedUser, err
}

func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User with that ID not found")
	}

	return user, nil
}
