package repos

import "backend/lib/models"

type IResultRepo interface {
	Save(*models.SuiteResult) (int, error)
}
