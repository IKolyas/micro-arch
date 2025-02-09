package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/IKolyas/otus-highload/internal/domain"
	"github.com/IKolyas/otus-highload/internal/infrastructure/database"
)

type UserRepository struct {
	Connection *database.Connection
}

func (r *UserRepository) GetAuthData(login string) (*domain.User, error) {
	pull := r.Connection.QueryFromReplica()
	if pull == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	row := "SELECT id, login, password FROM users WHERE login = $1"

	err := pull.QueryRow(context.Background(), row, login).Scan(&user.ID, &user.Login, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	pull := r.Connection.QueryFromReplica()
	if pull == nil {
		return nil, errors.New("database connection is nil")
	}

	user := domain.User{}

	row := "SELECT id, login, first_name, second_name, gender, birthdate, biography, city FROM users WHERE id = $1"

	err := pull.QueryRow(context.Background(), row, id).Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Gender, &user.Birthdate, &user.Biography, &user.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Find(fields map[string]string) ([]domain.User, error) {
	pull := r.Connection.QueryFromReplica()
	if pull == nil {
		return nil, errors.New("database connection is nil")
	}

	// Build the dynamic query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, login, first_name, second_name, gender, birthdate, biography, city FROM users WHERE ")

	var conditions []string
	var args []any

	count := 1
	for field, value := range fields {
		conditions = append(conditions, field+" LIKE $"+strconv.Itoa(count))
		args = append(args, value+"%")
		count++
	}

	queryBuilder.WriteString(strings.Join(conditions, " AND "))
	queryBuilder.WriteString(" ORDER BY id")

	query := queryBuilder.String()

	rows, err := pull.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Login, &user.FirstName, &user.SecondName, &user.Gender, &user.Birthdate, &user.Biography, &user.City); err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Save(user *domain.User) (res int, err error) {
	pull := r.Connection.QueryToMaster()
	if pull == nil {
		return 0, errors.New("database connection is err")
	}

	row := "INSERT INTO users (login, password, first_name, second_name, gender, birthdate, biography, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	lastInsert := 0
	err = pull.QueryRow(context.Background(), row, user.Login, user.Password, user.FirstName, user.SecondName, &user.Gender, user.Birthdate, user.Biography, user.City).Scan(&lastInsert)
	if err != nil {
		return 0, errors.New("DB error: " + err.Error())
	}

	user.ID = lastInsert

	return lastInsert, nil
}
