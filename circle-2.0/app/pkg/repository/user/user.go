package user

import (
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"errors"

	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

type UserRepo interface {
	CreateUser(tx *gorm.DB, user model.Users) error
	GetAllUsers() ([]model.Users, error)
	GetUserByID(id string) (model.Users, error)
	GetUserByName(name string) (model.Users, error)
	GetUserByEmail(email string) (model.Users, error)
	UpdateUser(id string, user map[string]any) error
	DeleteUser(id string) error
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		DB: db,
	}
}

func (u *userRepo) CreateUser(tx *gorm.DB, user model.Users) error {
	if tx != nil {
		return tx.Create(&user).Error
	} else {
		return u.DB.Create(&user).Error
	}
}

func (u *userRepo) GetAllUsers() ([]model.Users, error) {
	var users []model.Users
	err := u.DB.Find(&users).Error
	return users, err
}

func (u *userRepo) GetUserByID(id string) (model.Users, error) {
	var user model.Users
	err := u.DB.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New(customError.NotFoundError("user"))
	}

	return user, err
}

func (u *userRepo) GetUserByName(name string) (model.Users, error) {
	var user model.Users
	err := u.DB.Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New(customError.NotFoundError("user"))
	}

	return user, err
}

func (u *userRepo) GetUserByEmail(email string) (model.Users, error) {
	var user model.Users
	err := u.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New(customError.NotFoundError("user"))
	}
	return user, err
}

func (u *userRepo) UpdateUser(id string, user map[string]any) error {
	updateActivity := u.DB.Model(&model.Users{}).Where("id = ?", id).Updates(user)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("user"))
	}

	return nil
}

func (u *userRepo) DeleteUser(id string) error {
	// You can use this form if the parameter is not primary key
	// return u.DB.Where("id = ?", id).Delete(&model.User{}).Error

	deleteActivity := u.DB.Delete(&model.Users{}, "id = ?", id)
	if err := deleteActivity.Error; err != nil {
		return err
	}

	if deleteActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("user"))

	}

	return nil
}
