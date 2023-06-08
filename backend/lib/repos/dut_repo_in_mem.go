package repos

import (
	"backend/lib/models"
	"backend/utils"
)

type DUTRepoInMem struct {
	next_id int
	data    map[int]*models.DUT
}

func NewDUTRepoInMem() IDUTRepo {
	return &DUTRepoInMem{
		next_id: 0,
		data:    map[int]*models.DUT{},
	}
}

// Save save DUT information.
func (d *DUTRepoInMem) Save(model, board string) (int, error) {
	d.next_id += 1
	d.data[d.next_id] = &models.DUT{
		Id:    d.next_id,
		Model: model,
		Board: board,
	}
	return d.next_id, nil
}

// Get get the DUT by ID
func (d *DUTRepoInMem) Get(id int) (*models.DUT, error) {
	dut, found := d.data[id]
	if found {
		return dut, nil
	}

	return nil, utils.NotFound
}

// Get the DUT's ID by model and board
func (d *DUTRepoInMem) GetIdBy(model, board string) (int, error) {
	for id, item := range d.data {
		if item.Model == model && item.Board == board {
			return id, nil
		}
	}

	return -1, utils.NotFound
}

// SaveIfNotExit save DUT information if the DUT is not
// exist. Otherwise return the DUT's ID
func (d *DUTRepoInMem) SaveIfNotExist(model, board string) (int, error) {
  id, err := d.GetIdBy(model, board)
  if err == utils.NotFound {
    return d.Save(model, board)
  }
  return id, err
}
