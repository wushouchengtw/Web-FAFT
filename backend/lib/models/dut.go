package models

type DUT struct {
	Id    string `db:"dut_id"`
	Model string `db:"model"`
	Board string `db:"board"`
}
