package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// interface Service
type Service interface {
	RegisterUser(input RegisterInputUser) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
}

// Struct service
type service struct {
	repository Repository
}

// func new service
func NewService(repository Repository) *service {
	return &service{repository}
}

// func registerUser
func (s *service) RegisterUser(input RegisterInputUser) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	// enkripsi password yang diinput menggunakan encrypt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	// convert password yang sudah di enkripsi kedalam passwordHash yang ada pada struct user
	user.PasswordHash = string(passwordHash)
	user.Role = "User"

	// menjalankan service -> repository ->save(data user) untuk menginputkan data ke db
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, err
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("tidak ada user yang menggunakan email itu")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil

}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, err
}

// kegunaan service
// mapping struct input ke struct user
// simpan struct user melalui repository
