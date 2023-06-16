package dut

import (
	"backend/lib/models"
	"backend/utils"
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

// To-do
func (dm *DutMem) GetDUTCache() {
}

func (dm *DutMem) SaveIfNotExist(board, model string) (int, error) {
	return 0, nil
}
func (dm *DutMem) GetIdBy(board, model string) (int, error) {
	return 0, nil
}

func (dm *DutMem) FlahsDUTCache(id int, board, model string) {

}
