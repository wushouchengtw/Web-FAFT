package test

type Itest interface {
	SaveIfNotExist(testName string) error
	GetIdFromDBBy(testName string) (*string, error)
	SaveDB(id, testName string) error

	GetCache()
	GetIdByCache(testName string) (*string, error)
	FlashCache(id, testName string)
}
