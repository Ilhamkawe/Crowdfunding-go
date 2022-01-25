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
	GetUserByID(ID int) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(input FormUpdateRegister) (User, error)
	UpdateUserInfo(ID int, input UpdateInfoUserInput) (User, error)
	ChangePassword(input ChangePasswordInput) (User, error)
}

// Struct service
type service struct {
	repository Repository
}

// mendeklarasikan service untuk user
func NewService(repository Repository) *service {
	return &service{repository}
}

// fungsi untuk melakukan register user
func (s *service) RegisterUser(input RegisterInputUser) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.AvatarFileName = "images/default-user.jpg"
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

// func untuk proses login user
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

// memeriksa apakah email sudah dipakai atau belum
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

// meng update avatar dari user
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

func (s *service) UpdateUserInfo(ID int, input UpdateInfoUserInput) (User, error) {
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Occupation = input.Occupation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, err
}

// mengambiil data user berdasarkan UserID
func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindById(ID)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User tidak ditemukan")
	}

	return user, nil
}

// mengambil semua data user
func (s *service) GetAllUsers() ([]User, error) {
	users, err := s.repository.FindAll()

	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *service) UpdateUser(input FormUpdateRegister) (User, error) {
	user, err := s.repository.FindById(input.ID)

	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	updatedUser, err := s.repository.Update(user)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, err

}

func (s *service) ChangePassword(input ChangePasswordInput) (User, error) {

	user, err := s.repository.FindById(input.ID)

	if err != nil {
		return user, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)

	newPassword, err := s.repository.ChangePassword(user)
	if err != nil {
		return newPassword, err
	}
	return newPassword, nil
}

// kegunaan service
// mapping struct input ke struct user
// simpan struct user melalui repository
