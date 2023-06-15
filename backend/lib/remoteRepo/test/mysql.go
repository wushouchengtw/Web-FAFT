package test

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TestMySQL struct {
	db *sqlx.DB
}

func NewTestRepoInMySQL(db *sqlx.DB) Itest {
	return &TestMySQL{
		db: db,
	}
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
