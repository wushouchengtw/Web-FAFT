package test

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TestMySQL struct {
	db    *sqlx.DB
	mapId int
	cache map[int]*models.Test
}

func NewTestRepoInMySQL(db *sqlx.DB) Itest {
	return &TestMySQL{
		db:    db,
		mapId: 0,
		cache: map[int]*models.Test{},
	}
}

func (t *TestMySQL) GetCache() {
	dutList := []models.Test{}
	t.db.Select(&dutList, "SELECT * FROM Test")
	for _, dut := range dutList {
		t.mapId++
		t.cache[t.mapId] = &models.Test{Id: dut.Id, Name: dut.Name}
	}
}

func (t *TestMySQL) GetIdByCache(testName string) (int, error) {
	for _, dut := range t.cache {
		if dut.Name == testName {
			return dut.Id, nil
		}
	}
	return -1, utils.ErrNotFound
}

func (t *TestMySQL) FlashCache(id int, testName string) {
	t.mapId++
	t.cache[t.mapId] = &models.Test{Id: id, Name: testName}
}

func (t *TestMySQL) GetIdBy(testName string) (int, error) {
	var id int
	if err := t.db.Get(&id, "SELECT test_id FROM Test WHERE name=?", testName); err != nil {
		return -1, fmt.Errorf("failed to find %q in DB: %v", testName, err)
	}
	return id, nil
}

func (t *TestMySQL) Save(testName string) (int, error) {
	_, err := t.db.NamedExec("INSERT INTO Test(name) VALUES(:name)", []map[string]interface{}{
		{"name": testName},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert data into DB: %v", err)
	}
	return 0, nil
}

func (t *TestMySQL) SaveIfNotExist(testName string) (int, error) {
	test_id, err := t.GetIdByCache(testName)
	if err == utils.ErrNotFound {
		_, errSave := t.Save(testName)
		if errSave != nil {
			return -1, fmt.Errorf("failed to insert %q to DB: %v", testName, err)
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
