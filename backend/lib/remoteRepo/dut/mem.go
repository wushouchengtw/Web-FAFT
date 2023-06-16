package dut

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"
)

type DutMem struct {
	id   int
	data map[int]*models.DUT
}

func NewDutMem() IDUT {
	return &DutMem{
		id:   0,
		data: map[int]*models.DUT{},
	}
}

func (dm *DutMem) Save(board, model string) (int, error) {
	dm.id += 1
	dm.data[dm.id] = &models.DUT{
		Id:    dm.id,
		Board: board,
		Model: model,
	}
	return dm.id, nil
}

func (dm *DutMem) GetIdByCache(board, model string) (int, error) {
	for id, dut := range dm.data {
		if dut.Board == board && dut.Model == model {
			return id, nil
		}
	}
	return -1, utils.ErrNotFound
}

func (dm *DutMem) SaveIfNotExist(board, model string) (int, error) {
	dut_id, err := dm.GetIdBy(board, model)
	if err == utils.ErrNotFound {
		_, errSave := dm.Save(board, model)
		if errSave != nil {
			return -1, fmt.Errorf("failed to save (%q,%q) in DB: %v", board, model, err)
		} else {
			dut_id, errGetID := dm.GetIdBy(board, model)
			if errGetID != nil {
				return -1, fmt.Errorf("failed to find (%q,%q) in DB: %v", board, model, err)
			}
			dm.FlashCache(dut_id, board, model)
			return dut_id, nil
		}
	}
	return dut_id, nil
}

func (dm *DutMem) GetCache() {

}
func (dm *DutMem) GetIdBy(board, model string) (int, error) {
	return 0, nil
}

func (dm *DutMem) FlashCache(id int, board, model string) {

}
