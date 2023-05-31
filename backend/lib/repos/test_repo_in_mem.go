package repos

import (
	"backend/lib/models"
	"backend/utils"
)

type TestRepoInMem struct {
	next_id int
	data    map[int]*models.Test
}

func NewTestRepoInMem() ITestRepo {
	return &TestRepoInMem{
		next_id: 0,
		data:    map[int]*models.Test{},
	}
}

func (t *TestRepoInMem) Save(name string) (int, error) {
	t.next_id += 1
	t.data[t.next_id] = &models.Test{Name: name}
	return t.next_id, nil
}

func (t *TestRepoInMem) Get(id int) (*models.Test, error) {
	test, found := t.data[id]
	if found {
		return test, nil
	}

	return nil, utils.NotFound
}

func (t *TestRepoInMem) GetIdBy(name string) (int, error) {
	for id, item := range t.data {
		if item.Name == name {
			return id, nil
		}
	}

	return -1, utils.NotFound
}
