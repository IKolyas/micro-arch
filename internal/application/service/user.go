package service

import (
	"errors"

	"github.com/IKolyas/otus-highload/internal/application/utils"
	"github.com/IKolyas/otus-highload/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type User struct{}

func (u *User) Register(user *domain.User, r domain.Repository[domain.User]) (userId int, err error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, errors.New("Generate password error")
	}

	user.Password = string(passwd)

	userId, err = r.Save(user)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (u *User) Login(login string, password string, r domain.UserRepository) (*domain.User, error) {

	user, err := r.GetAuthData(login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	res := utils.ComparePassword(user.Password, password)

	if !res {
		return nil, errors.New("incorrect password")
	}

	return user, nil
}

func (u *User) GetByID(id int, r domain.Repository[domain.User]) (map[string]interface{}, error) {

	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return map[string]interface{}{
		"id":         user.ID,
		"login":      user.Login,
		"firstName":  user.FirstName,
		"secondName": user.SecondName,
		"gender":     user.Gender,
		"birthdate":  user.Birthdate,
		"biography":  user.Biography,
		"city":       user.City,
	}, nil
}

func (u *User) Find(fields map[string]string, r domain.Repository[domain.User]) ([]any, error) {

	users, err := r.Find(fields)
	if err != nil {
		return nil, err
	}

	var res []any
	for _, user := range users {
		res = append(res, map[string]any{
			"id":         user.ID,
			"login":      user.Login,
			"firstName":  user.FirstName,
			"secondName": user.SecondName,
			"gender":     user.Gender,
			"birthdate":  user.Birthdate,
			"biography":  user.Biography,
			"city":       user.City,
		})
	}

	return res, nil
}
