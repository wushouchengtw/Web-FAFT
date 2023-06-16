package dut

type IDUT interface {
	SaveIfNotExist(board, model string) (int, error)
	Save(board, model string) (int, error)
	GetIdBy(board, model string) (int, error)

	GetCache()
	GetIdByCache(board, name string) (int, error)
	FlashCache(id int, board, model string)
}
