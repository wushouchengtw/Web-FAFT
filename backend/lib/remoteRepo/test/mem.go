package test

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/google/uuid"
)

type TestMem struct {
	data map[string]*models.Test
}

func NewTestRepoInMem() Itest {
	return &TestMem{
		data: map[string]*models.Test{},
	}
}

func (t *TestMem) GetIdFromDBBy(testName string) (*string, error) {
	for id, data := range t.data {
		if data.Name == testName {
			return &id, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (t *TestMem) SaveDB(id, testName string) error {
	t.data[id] = &models.Test{Id: id, Name: testName}
	return nil
}

func (t *TestMem) SaveIfNotExist(testName string) (*string, error) {
	testID, err := t.GetIdFromDBBy(testName)
	if err != nil {
		id := uuid.New().String()
		if err := t.SaveDB(id, testName); err != nil {
			return nil, fmt.Errorf("failed to ")
		}
	}
	return testID, nil
}

func (t *TestMem) GetCache()

func (t *TestMem) GetIdByCache(testName string) (*string, error)

func (t *TestMem) FlashCache(id, testName string)
