package store

import (
	"github.com/amannlalwani/Notes-app-using-gofr/models"
	"gofr.dev/pkg/gofr"
)

type Store interface {
	Get(ctx *gofr.Context) ([]models.Note, error)
	Create(ctx *gofr.Context, note models.Note) (models.Note, error)
	Update(ctx *gofr.Context, id int, note models.Note) (models.Note, error)
	Delete(ctx *gofr.Context, id int) (models.Note, error)
}
