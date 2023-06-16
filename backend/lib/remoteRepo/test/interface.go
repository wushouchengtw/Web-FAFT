package test

type Itest interface {
	SaveIfNotExist(testName string) (int, error)
	GetIdBy(testName string) (int, error)
	Save(testName string) (int, error)

	GetCache()
	GetIdByCache(testName string) (int, error)
	FlashCache(id int, testName string)
}
