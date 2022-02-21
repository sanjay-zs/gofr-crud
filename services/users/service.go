package users

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"gofr-crud/models"
	"gofr-crud/stores"
	"strconv"
)

type Service struct {
	store stores.Users
}

func New(users stores.Users) Service {
	return Service{store: users}
}

// Create inserts new user into the database
func (s Service) Create(ctx *gofr.Context, user *models.User) error {
	err := s.store.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetAll fetches all users from the database
func (s Service) GetAll(ctx *gofr.Context) ([]models.User, error) {
	users, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID fetches user details for given id
func (s Service) GetByID(ctx *gofr.Context, id int) (models.User, error) {
	if id <= 0 {
		return models.User{}, errors.InvalidParam{Param: []string{"id"}}
	}

	user, err := s.store.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Update updates the existing user details
func (s Service) Update(ctx *gofr.Context, user *models.User) error {
	id := user.ID

	// validating whether id is valid id or not
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	// checking whether the entity exists or not
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		return errors.EntityNotFound{
			Entity: "user",
			ID:     strconv.Itoa(id),
		}
	}

	err = s.store.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes the user with given id from database
func (s Service) Delete(ctx *gofr.Context, id int) error {
	// validating whether id is valid id or not
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	// checking whether the entity exists or not
	_, err := s.store.GetByID(ctx, id)
	if err != nil {
		return errors.EntityNotFound{
			Entity: "user",
			ID:     strconv.Itoa(id),
		}
	}

	err = s.store.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
