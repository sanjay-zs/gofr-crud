package users

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"gofr-crud/models"
)

type User struct{}

func New() User {
	return User{}
}

// Create inserts a new user in the database
func (u User) Create(ctx *gofr.Context, user *models.User) error {
	_, err := ctx.DB().ExecContext(ctx, insertUser, user.ID, user.Name, user.Age)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// GetAll fetches all the users from the database
func (u User) GetAll(ctx *gofr.Context) ([]models.User, error) {
	rows, err := ctx.DB().QueryContext(ctx, getUsers)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, errors.Error("Scan error")
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByID fetches the user details for given id
func (u User) GetByID(ctx *gofr.Context, id int) (models.User, error) {
	var user models.User
	err := ctx.DB().QueryRowContext(ctx, getUserByID, id).Scan(&user.ID, &user.Name, &user.Age)

	if err != nil {
		return models.User{}, errors.DB{Err: err}
	}

	return user, nil
}

// Update updates the user details
func (u User) Update(ctx *gofr.Context, user *models.User) error {
	_, err := ctx.DB().ExecContext(ctx, updateUser, user.Name, user.Age, user.ID)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

// Delete removes the user details from database
func (u User) Delete(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, deleteUser, id)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}
