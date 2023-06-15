package models 

type DUT struct {
				Id    int   		`db:"dut_id"`
				Model string		`db:"model"`
				Board string		`db:"board"`
}
