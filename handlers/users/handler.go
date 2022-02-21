package users

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"gofr-crud/models"
	"gofr-crud/services"
	"strconv"
)

type Handler struct {
	service services.Users
}

func New(service services.Users) Handler {
	return Handler{service: service}
}

// Create takes the clients request to create entity in database
func (h Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var user models.User

	var result models.Response

	if err := ctx.Bind(&user); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	err := h.service.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	result = models.Response{
		User:       user,
		Message:    "Inserted user successfully",
		StatusCode: 200,
	}

	return result, nil
}

// GetAll writes all the users from the database
func (h Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.service.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	res := models.Response{
		User:       resp,
		Message:    "Retrieved User Successfully",
		StatusCode: 200,
	}

	return res, nil
}

// GetByID writes the response based on ID of the resp
func (h Handler) GetByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := h.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := models.Response{
		User:       resp,
		Message:    "Retrieved User Successfully",
		StatusCode: 200,
	}

	return res, nil
}

// Update writes the updated resp entity in the database
func (h Handler) Update(ctx *gofr.Context) (interface{}, error) {
	var user models.User

	var result models.Response

	if err := ctx.Bind(&user); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	err := h.service.Update(ctx, &user)

	if err != nil {
		return nil, err
	}

	result = models.Response{
		User:       user,
		Message:    "Updated user successfully",
		StatusCode: 200,
	}

	return result, nil
}

// Delete removes the resp from database based on ID
func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	res := models.Response{
		User:       nil,
		Message:    "Deleted User Successfully",
		StatusCode: 200,
	}

	return res, nil
}
