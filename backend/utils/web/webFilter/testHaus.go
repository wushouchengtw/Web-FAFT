package webFilter

import (
	webmodels "backend/utils/web/web_models"
	"reflect"
)

type TestHaus struct {
	TableName string
}

func NewTestHausSearch() QueryMethod {
	return &TestHaus{
		TableName: "Result",
	}
}

func (t *TestHaus) ToConditions(params *webmodels.QueryParameter) ([]webmodels.QueryValue, error) {
	output := []webmodels.QueryValue{}

	rv := reflect.ValueOf(&params)
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Interface() != "" {
			fieldName := rt.Field(i).Name
			switch fieldName {
			case "StartDate":
				output = append(output, webmodels.QueryValue{Filter: fieldName + " >= ?", Value: rv.Field(i).String(), Condition: "AND"})
			case "EndDate":
				output = append(output, webmodels.QueryValue{Filter: fieldName + " <= ?", Value: rv.Field(i).String(), Condition: "AND"})
			default:
				output = append(output, webmodels.QueryValue{Filter: fieldName + " = ?", Value: rv.Field(i).String(), Condition: "AND"})
			}
		}
	}
	return output, nil
}
