package dut

type IDUT interface {
	Save(board, model string) (int, error)
	GetIdBy(board, name string) (int, error)
	GetDUTCache()
}
