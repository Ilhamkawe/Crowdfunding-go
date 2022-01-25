package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindById(ID int) (User, error)
	Update(user User) (User, error)
	FindAll() ([]User, error)
	ChangePassword(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

// mendeklarasikan repository untuk user
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// menyimpan user baru
func (r *repository) Save(user User) (User, error) {
	// insert into user values (user input)
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// mencari user berdasarkan email
func (r *repository) FindByEmail(email string) (User, error) {
	var user User

	// select * from user where email = ?
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

// mencari user berdasarkan id
func (r *repository) FindById(ID int) (User, error) {
	var user User

	// select * from user where id = ?
	err := r.db.Where("id = ?", ID).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// mengupdate data user
func (r *repository) Update(user User) (User, error) {
	// update user set user where id = user.id
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// mengambil semua data user
func (r *repository) FindAll() ([]User, error) {
	var user []User
	// select * from user
	err := r.db.Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// update password
func (r *repository) ChangePassword(user User) (User, error) {
	err := r.db.Model(&user).Updates(User{PasswordHash: user.PasswordHash}).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
