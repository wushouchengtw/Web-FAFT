package webFilter

import (
	webmodels "backend/utils/web/web_models"
)

type QueryMethod interface {
	ToConditions(params *webmodels.QueryParameter) ([]webmodels.QueryValue, error)
}
