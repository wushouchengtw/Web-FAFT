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

func (d *DutMem) Save(board, model string) (int, error) {
	d.id += 1
	d.data[d.id] = &models.DUT{
		Id:    d.id,
		Board: board,
		Model: model,
	}
	return d.id, nil
}

func (d *DutMem) GetIdBy(board, name string) (int, error) {
	for id, dut := range d.data {
		if dut.Board == board && dut.Model == dut.Model {
			return id, nil
		}
	}
	return -1, utils.ErrNotFound
}
