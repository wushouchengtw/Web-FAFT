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

// Save the test case
func (t *TestRepoInMem) Save(name string) (int, error) {
	t.next_id += 1
	t.data[t.next_id] = &models.Test{Name: name}
	return t.next_id, nil
}

// Get the test case by id
func (t *TestRepoInMem) Get(id int) (*models.Test, error) {
	test, found := t.data[id]
	if found {
		return test, nil
	}

	return nil, utils.NotFound
}

// GetIdBy get the test case by test case' name
func (t *TestRepoInMem) GetIdBy(name string) (int, error) {
	for id, item := range t.data {
		if item.Name == name {
			return id, nil
		}
	}

	return -1, utils.NotFound
}

// SaveIfNotExist save the test case if it is not exist
func (t* TestRepoInMem) SaveIfNotExist(name string) (int, error) {
  id, err := t.GetIdBy(name)
  if err == utils.NotFound {
    return t.Save(name)
  }
  return id, err
}
