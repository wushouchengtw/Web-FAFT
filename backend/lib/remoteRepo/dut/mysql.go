package dut

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DUTMySQL struct {
	db    *sqlx.DB
	cache map[string]*models.DUT
}

func NewDUTRepoInMySQL(db *sqlx.DB) IDUT {
	return &DUTMySQL{
		db:    db,
		cache: map[string]*models.DUT{},
	}
}

func (ds *DUTMySQL) GetCache() {
	dutList := []models.DUT{}
	ds.db.Select(&dutList, "SELECT * FROM DUT")
	for _, dut := range dutList {
		ds.cache[dut.Id] = &models.DUT{Id: dut.Id, Board: dut.Board, Model: dut.Model}
	}
}

func (ds *DUTMySQL) GetIdByCache(board, model string) (*string, error) {
	for _, dut := range ds.cache {
		if dut.Board == board && dut.Model == model {
			return &dut.Id, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (ds *DUTMySQL) FlashCache(id, board, model string) {
	ds.cache[id] = &models.DUT{Id: id, Board: board, Model: model}
}

func (ds *DUTMySQL) SaveDB(id, board, model string) error {
	_, err := ds.db.NamedExec("INSERT INTO DUT(dut_id,board, model) values(:dut_id,:board,:model)", []map[string]interface{}{
		{
			"dut_id": id,
			"board":  board,
			"modle":  model,
		}})
	if err != nil {
		return fmt.Errorf("failed to insert data into DB: %v", err)
	}
	return nil
}

func (ds *DUTMySQL) GetIdFromDBBy(board, model string) (*string, error) {
	var id string
	if err := ds.db.Get(&id, "SELECT dut_id FROM DUT WHERE board = ? AND model = ?", board, model); err != nil {
		return nil, fmt.Errorf("failed to find (%q,%q) in DB: %v", board, model, err)
	}
	return &id, nil
}

func (ds *DUTMySQL) SaveIfNotExist(board, model string) (*string, error) {
	// Search cache
	dutID, err := ds.GetIdByCache(board, model)
	// Cache miss
	if err == utils.ErrNotFound {
		// Search DB
		dut_id, err := ds.GetIdFromDBBy(board, model)
		// Not found in DB
		if err != nil {
			id := uuid.New().String()
			errSave := ds.SaveDB(id, board, model)
			if errSave != nil {
				return nil, fmt.Errorf("failed to save in DB: %v", err)
			}
			ds.FlashCache(id, board, model)
		} else {
			ds.FlashCache(*dut_id, board, model)
		}
		return dut_id, nil
	}
	return dutID, nil
}
