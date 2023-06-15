package test

import (
	"backend/lib/models"
	"backend/utils"
)

type TestMem struct {
	id   int
	data map[int]*models.Test
}

func NewTestRepoInMem() Itest {
	return &TestMem{
		id:   0,
		data: map[int]*models.Test{},
	}
}

func (t *TestMem) GetIdBy(testName string) (int, error) {
	for id, data := range t.data {
		if data.Name == testName {
			return id, nil
		}
	}
	return -1, utils.ErrNotFound
}

func (t *TestMem) Save(testName string) (int, error) {
	t.id += 1
	t.data[t.id] = &models.Test{Name: testName}
	return t.id, nil
}
