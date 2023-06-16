package dut

type IDUT interface {
	SaveIfNotExist(board, model string) (int, error)
	Save(board, model string) (int, error)
	GetIdBy(board, model string) (int, error)

	GetDUTCache()
	GetIdByCache(board, name string) (int, error)
	FlahsDUTCache(id int, board, model string)
}
