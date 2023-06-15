package test

type Itest interface {
	GetIdBy(testName string) (int, error)
	Save(testName string) (int, error)
}
