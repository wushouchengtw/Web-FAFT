package dut

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DUTMySQL struct {
	db   *sqlx.DB
	data map[int]*models.DUT
}

func NewDUTRepoInMySQL(db *sqlx.DB) IDUT {
	return &DUTMySQL{
		db:   db,
		data: map[int]*models.DUT{},
	}
}

func (ds *DUTMySQL) GetIdBy(board, model string) (int, error) {
	//var id int
	//if err := t.db.Get(&id, "SELECT dut_id FROM DUT WHERE board = ? AND model = ?", board, model); err != nil {
	//	return -1, fmt.Errorf("failed to find (%q,%q) in DB: %v", board, model, err)
	//}
	//return id, nil
	for id, dut := range ds.data {
		if dut.Board == board && dut.Model == model {
			return id, nil
		}
	}
	return -1, utils.ErrNotFound
}

func SaveIfNotExist() {

}

func (ds *DUTMySQL) Save(board, model string) (int, error) {
	_, err := ds.db.NamedExec("INSERT INTO DUT(board, model) values(:board,:model)", []map[string]interface{}{
		{
			"board": board,
			"modle": model,
		}})
	if err != nil {
		return -1, fmt.Errorf("failed to insert data into DB: %v", err)
	}
	return 0, nil
}

func (ds *DUTMySQL) GetDUTCache() {
	dutList := []models.DUT{}
	ds.db.Select(&dutList, "SELECT * FROM DUT")
	for _, dut := range dutList {
		ds.data[dut.Id] = &models.DUT{Id: dut.Id, Board: dut.Board, Model: dut.Model}
	}
}
