package repos

import "backend/lib/models"

type IDUTRepo interface {
	Save(model, board string) (int, error)

	Get(int) (*models.DUT, error)

	GetIdBy(model, board string) (int, error)
}
