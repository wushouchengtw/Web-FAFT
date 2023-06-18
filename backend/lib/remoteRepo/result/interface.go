package result

import "backend/lib/models"

type IResut interface {
	Save(v *models.Result) (int, error)
	SearchTestHaus(search models.SearchParams)
}
