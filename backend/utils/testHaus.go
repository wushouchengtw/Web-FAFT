package utils

import (
	"fmt"
	"reflect"
	"time"
)

type QueryParameter struct {
	StartDate time.Time `field:"StartDate" op:" >= ?"`
	EndDate   time.Time `field:"EndDate" op:" <= ?"`
	Board     string    `field:"Board" op:" = ?"`
	Reason    string    `field:"Reason" op:"none" `
	Name      string    `field:"Name" op:" = ?"`
	Status    bool      `field:"Status" op:" = ?"`
}

type QueryOperation struct {
	Where string
	Value interface{}
}

// To-do : Clean up the code
func ToConditions(q interface{}) ([]QueryOperation, error) {
	t := reflect.TypeOf(q)
	v := reflect.ValueOf(q)

	values := make([]QueryOperation, 0, t.NumField())

	for idx := 0; idx < t.NumField(); idx++ {
		field := t.Field(idx)
		val := v.Field(idx)

		if val.IsValid() && !val.IsZero() {
			v := val.Interface()
			q := QueryOperation{}

			switch field.Type.Kind() {
			case reflect.Struct:
				if field.Type == reflect.TypeOf(time.Time{}) {
					q.Where = fmt.Sprintf("%s%s", field.Tag.Get("field"), field.Tag.Get("op"))
					q.Value = v.(time.Time)
				}
			case reflect.String:
				if field.Tag.Get("op") != "none" {
					q.Where = fmt.Sprintf("%s%s", field.Tag.Get("field"), field.Tag.Get("op"))
					q.Value = v.(string)
				} else {
					q.Where = v.(string)
				}
			case reflect.Bool:
				q.Where = fmt.Sprintf("%s%s", field.Tag.Get("field"), field.Tag.Get("op"))
				q.Value = v.(bool)
			default:
				fmt.Printf("Un-supported type {%v}", field.Type.Kind())
			}
			values = append(values, q)
		}
	}
	return values, nil
}
