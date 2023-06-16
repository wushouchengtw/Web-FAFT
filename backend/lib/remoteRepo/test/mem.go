package test

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"
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

func (t *TestMem) SaveIfNotExist(testName string) (int, error) {
	test_id, err := t.GetIdBy(testName)
	if err == utils.ErrNotFound {
		_, errSave := t.Save(testName)
		if errSave != nil {
			return -1, fmt.Errorf("failed to save %q in DB: %v", testName, err)
		} else {
			test_id, errGetID := t.GetIdBy(testName)
			if errGetID != nil {
				return -1, fmt.Errorf("failed to find %q in DB: %v", testName, err)
			}
			t.FlashCache(test_id, testName)
			return test_id, nil
		}
	}
	return test_id, nil
}

func (t *TestMem) GetCache() {

}

func (t *TestMem) GetIdByCache(testName string) (int, error) {
	return 0, nil
}

func (t *TestMem) FlashCache(id int, testName string) {

}
