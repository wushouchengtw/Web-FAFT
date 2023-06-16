package dut

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/google/uuid"
)

type DutMem struct {
	data map[string]*models.DUT
}

func NewDutMem() IDUT {
	return &DutMem{
		data: map[string]*models.DUT{},
	}
}

func (dm *DutMem) SaveDB(id, board, model string) error {
	dm.data[id] = &models.DUT{
		Id:    id,
		Board: board,
		Model: model,
	}
	return nil
}

func (dm *DutMem) GetIdFromDBBy(board, model string) (*string, error) {
	for id, dut := range dm.data {
		if dut.Board == board && dut.Model == model {
			return &id, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (dm *DutMem) SaveIfNotExist(board, model string) error {
	_, err := dm.GetIdFromDBBy(board, model)
	if err == utils.ErrNotFound {
		id := uuid.New().String()
		errSave := dm.SaveDB(id, board, model)
		if errSave != nil {
			return fmt.Errorf("failed to insert (%q,%q) to DB: %v", board, model, err)
		}
	}
	return nil
}

func (dm *DutMem) GetCache()

func (dm *DutMem) GetIdByCache(board, model string) (*string, error)

func (dm *DutMem) FlashCache(id, board, model string)
