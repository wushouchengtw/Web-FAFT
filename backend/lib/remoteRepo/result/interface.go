package result

import (
	"backend/lib/models"
	"backend/utils"
)

type IResut interface {
	Save(v *models.Result) (int, error)
	SearchTestHaus(params utils.QueryParameter) ([]models.RawDataFromResult, error)
}
