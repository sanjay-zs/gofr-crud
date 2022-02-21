package stores

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"gofr-crud/models"
)

type Users interface {
	Create(ctx *gofr.Context, user *models.User) error
	GetAll(ctx *gofr.Context) ([]models.User, error)
	GetByID(ctx *gofr.Context, id int) (models.User, error)
	Update(ctx *gofr.Context, user *models.User) error
	Delete(ctx *gofr.Context, id int) error
}
