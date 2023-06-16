package dut

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DUTMySQL struct {
	db    *sqlx.DB
	mapId int
	cache map[int]*models.DUT
}

func NewDUTRepoInMySQL(db *sqlx.DB) IDUT {
	return &DUTMySQL{
		db:    db,
		mapId: 0,
		cache: map[int]*models.DUT{},
	}
}

func (ds *DUTMySQL) GetCache() {
	dutList := []models.DUT{}
	ds.db.Select(&dutList, "SELECT * FROM DUT")
	for _, dut := range dutList {
		ds.mapId++
		ds.cache[ds.mapId] = &models.DUT{Id: dut.Id, Board: dut.Board, Model: dut.Model}
	}
}

func (ds *DUTMySQL) GetIdByCache(board, model string) (int, error) {
	for _, dut := range ds.cache {
		if dut.Board == board && dut.Model == model {
			return dut.Id, nil
		}
	}
	return -1, utils.ErrNotFound
}

func (ds *DUTMySQL) FlashCache(id int, board, model string) {
	ds.mapId++
	ds.cache[ds.mapId] = &models.DUT{Id: id, Board: board, Model: model}
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

func (ds *DUTMySQL) GetIdBy(board, model string) (int, error) {
	var id int
	if err := ds.db.Get(&id, "SELECT dut_id FROM DUT WHERE board = ? AND model = ?", board, model); err != nil {
		return -1, fmt.Errorf("failed to find (%q,%q) in DB: %v", board, model, err)
	}
	return id, nil
}

func (ds *DUTMySQL) SaveIfNotExist(board, model string) (int, error) {
	dut_id, err := ds.GetIdByCache(board, model)
	if err == utils.ErrNotFound {
		_, errSave := ds.Save(board, model)
		if errSave != nil {
			return -1, fmt.Errorf("failed to insert (%s,%s) to DB: %v", board, model, err)
		} else {
			dut_id, errGetID := ds.GetIdBy(board, model)
			if errGetID != nil {
				return -1, fmt.Errorf("failed to find (%s,%s) in DB: %v", board, model, err)
			}
			ds.FlashCache(dut_id, board, model)
			return dut_id, nil
		}
	}
	return dut_id, nil
}
