package repos

import "backend/lib/models"

type IDUTRepo interface {
  // Save save DUT information.
	Save(model, board string) (int, error)

  // Get get the DUT by ID
	Get(int) (*models.DUT, error)

  // Get the DUT's ID by model and board
	GetIdBy(model, board string) (int, error)

  // SaveIfNotExit save DUT information if the DUT is not
  // exist. Otherwise return the DUT's ID
  SaveIfNotExist(model, board string) (int, error)
}
