package controller

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IKolyas/otus-highload/internal/application/service"
	"github.com/IKolyas/otus-highload/internal/domain"
	"github.com/IKolyas/otus-highload/internal/infrastructure/actions"
	"github.com/IKolyas/otus-highload/internal/infrastructure/database"
	"github.com/IKolyas/otus-highload/internal/infrastructure/repository"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var xv = &XValidator{
	validator: validator.New(),
}

func (v XValidator) Validate(data interface{}) error {
	validationErrors := []ErrorResponse{}

	errs := xv.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	if len(validationErrors) > 0 {
		errMsgs := make([]string, 0)

		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | need to be realized '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	return nil
}

func Login(c *fiber.Ctx) error {
	u := new(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(u); err != nil {
		return err
	}

	user := struct {
		Login    string `validate:"required"`
		Password string `validate:"required"`
	}{Login: u.Login, Password: u.Password}

	err := xv.Validate(user)
	if err != nil {
		return err
	}
	uService := service.User{}
	r := repository.UserRepository{Connection: &database.DB}
	_, err = uService.Login(u.Login, u.Password, &r)
	if err != nil {
		return err
	}

	claims := jwt.MapClaims{
		"name":  u.Login,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Register(c *fiber.Ctx) error {
	body := new(struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		FirstName  string `json:"firstName"`
		SecondName string `json:"secondName"`
		Gender     string `json:"gender"`
		Birthdate  string `json:"birthdate"`
		Biography  string `json:"biography"`
		City       string `json:"city"`
	})

	if err := c.BodyParser(body); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	request := struct {
		Login      string `validate:"required,min=6"`
		Password   string `validate:"required,min=8"`
		FirstName  string `validate:"required,min=2"`
		SecondName string `validate:"required,min=2"`
		Gender     string `validate:"required"`
		Birthdate  string `validate:"required"`
		Biography  string `validate:"required"`
		City       string `validate:"required,min=2"`
	}{Login: body.Login, Password: body.Password, FirstName: body.FirstName, Gender: body.Gender, SecondName: body.SecondName, Birthdate: body.Birthdate, Biography: body.Biography, City: body.City}

	if err := xv.Validate(request); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	bd, err := time.Parse("01.02.2006", request.Birthdate)
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": errors.New("Invalid Birthdate field")})
	}

	uService := service.User{}
	user := domain.User{
		Login:      request.Login,
		Password:   request.Password,
		FirstName:  request.FirstName,
		SecondName: request.SecondName,
		Gender:     request.Gender,
		Birthdate:  bd,
		Biography:  request.Biography,
		City:       request.City,
	}

	id, err := uService.Register(&user, &repository.UserRepository{Connection: &database.DB})
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.JSON(fiber.Map{"status": "success", "ID": id})
}

// ------------------------- AUTH USER -------------------------

func GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": errors.New("Not parse id")})
	}

	uService := service.User{}

	profile, err := uService.GetByID(id, &repository.UserRepository{Connection: &database.DB})
	if err != nil {
		c.SendStatus(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"status": "error", "message": errors.New("Not found user")})
	}

	return c.JSON(fiber.Map{"status": "success", "data": profile})
}

func SearchUser(c *fiber.Ctx) error {
	fname := c.Params("firstName")
	lname := c.Params("secondName")

	fields := map[string]string{"first_name": fname, "second_name": lname}

	uService := service.User{}
	users, err := uService.Find(fields, &repository.UserRepository{Connection: &database.DB})
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.JSON(fiber.Map{"status": "success", "data": users})
}

func FakerUser(c *fiber.Ctx) error {
	count, err := strconv.Atoi(c.Params("count"))
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": err})
	}

	ok, errors := actions.CreateRandomUsers(count, &repository.UserRepository{Connection: &database.DB})

	return c.JSON(fiber.Map{"status": "success", "data": fiber.Map{
		"created": ok,
		"error":   errors,
	}})

}
