package repos

import "backend/lib/models"

type ITestRepo interface {
	Save(name string) (int, error)

	Get(int) (*models.Test, error)

	GetIdBy(name string) (int, error)
}
