package test

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TestMySQL struct {
	db    *sqlx.DB
	cache map[string]*models.Test
}

func NewTestRepoInMySQL(db *sqlx.DB) Itest {
	return &TestMySQL{
		db:    db,
		cache: map[string]*models.Test{},
	}
}

func (t *TestMySQL) GetCache() {
	testList := []models.Test{}
	t.db.Select(&testList, "SELECT * FROM Test")
	for _, test := range testList {
		t.cache[test.Id] = &models.Test{Id: test.Id, Name: test.Name}
	}
}

func (t *TestMySQL) GetIdByCache(testName string) (*string, error) {
	for _, test := range t.cache {
		if test.Name == testName {
			return &test.Id, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (t *TestMySQL) FlashCache(id, testName string) {
	t.cache[id] = &models.Test{Id: id, Name: testName}
}

func (t *TestMySQL) GetIdFromDBBy(testName string) (*string, error) {
	var id string
	if err := t.db.Get(&id, "SELECT test_id FROM Test WHERE name=?", testName); err != nil {
		return nil, fmt.Errorf("failed to find %q in DB: %v", testName, err)
	}
	return &id, nil
}

func (t *TestMySQL) SaveDB(id, testName string) error {
	_, err := t.db.NamedExec("INSERT INTO Test(test_id,name) VALUES(:test_id,:name)", []map[string]interface{}{
		{
			"id":   id,
			"name": testName,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to insert data into DB: %v", err)
	}
	return nil
}

func (t *TestMySQL) SaveIfNotExist(testName string) error {
	_, err := t.GetIdByCache(testName)
	if err == utils.ErrNotFound {
		test_id, err := t.GetIdFromDBBy(testName)
		if err != nil {
			if err := t.SaveDB(*test_id, testName); err != nil {
				return fmt.Errorf("failed to save %q in DB: %v", testName, err)
			}
		}
		t.FlashCache(*test_id, testName)
		return nil
	}
	return nil
}
