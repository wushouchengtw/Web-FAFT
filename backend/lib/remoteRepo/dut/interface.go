package dut

type IDUT interface {
	SaveIfNotExist(board, model string) (*string, error)
	SaveDB(id, board, model string) error
	GetIdFromDBBy(board, model string) (*string, error)

	GetCache()
	GetIdByCache(board, model string) (*string, error)
	FlashCache(id, board, model string)
}
