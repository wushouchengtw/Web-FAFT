package repos

import "backend/lib/models"

type ITestRepo interface {
  // Save the test case
	Save(name string) (int, error)

  // Get the test case by id
	Get(int) (*models.Test, error)

  // GetIdBy get the test case by test case' name
	GetIdBy(name string) (int, error)

  // SaveIfNotExist save the test case if it is not exist
  SaveIfNotExist(name string) (int, error)
}
