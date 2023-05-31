package repos

import (
	"backend/lib/models"
	"log"
)

type ResultRepoInMem struct {
	next_id int
	data    map[int]*models.SuiteResult
}

func NewResultRepoInMem() IResultRepo {
	return &ResultRepoInMem{
		next_id: 0,
		data:    map[int]*models.SuiteResult{},
	}
}

func (r *ResultRepoInMem) Save(value *models.SuiteResult) (int, error) {
	log.Printf("Save result: {%v}", value)
	r.next_id += 1
	r.data[r.next_id] = value
	return r.next_id, nil
}
