package result

import (
	"backend/lib/models"
	webmodels "backend/utils/web/web_models"
)

type IResut interface {
	Save(v *models.Result) (int, error)
	SearchTestHaus(params webmodels.QueryParameter)
}
