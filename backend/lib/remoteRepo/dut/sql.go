package dut

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DUTMySQL struct {
	db *sqlx.DB
}

func NewDUTRepoInMySQL(db *sqlx.DB) IDUT {
	return &DUTMySQL{
		db: db,
	}
}

func (t *DUTMySQL) GetIdBy(board, model string) (int, error) {
	var id int
	if err := t.db.Get(&id, "SELECT dut_id FROM DUT WHERE board = ? AND model = ?", board, model); err != nil {
		return -1, fmt.Errorf("failed to find (%q,%q) in DB: %v", board, model, err)
	}
	return id, nil
}

func (t *DUTMySQL) Save(board, model string) (int, error) {
	_, err := t.db.NamedExec("INSERT INTO DUT(board, model) values(:board,:model)", []map[string]interface{}{
		{
			"board": board,
			"modle": model,
		},
	})
	if err != nil {
		return -1, fmt.Errorf("failed to insert data into DB: %v", err)
	}
	return 0, nil
}
