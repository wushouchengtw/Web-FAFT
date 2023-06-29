package result

import (
	"backend/lib/models"
	"backend/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ResultMySQL struct {
	db *sqlx.DB
}

func NewResultRepoMySQL(db *sqlx.DB) IResut {
	return &ResultMySQL{
		db: db,
	}
}

func (r *ResultMySQL) Save(v *models.Result) (int, error) {
	_, err := r.db.NamedExec("INSERT INTO Result(result_id,time,duration,suite,dut_id,milestone,version,host,test_id,status,reason,firmware_RO_Version,firmware_RW_Version) VALUES(:result_id,:time,:duration,:suite,:dut_id,:milestone,:version,:host,:test_id,:status,:reason,:firmware_RO_Version,:firmware_RW_Version)", v)
	if err != nil {
		return -1, fmt.Errorf("failed to insert data [%v] into DB: ", v)
	}
	return 0, nil
}

// func (r *ResultMySQL) GetAllDUT() ([]models.DUT, error) {
// 	dut := []models.DUT{}
// 	if err := r.db.Select(dut, "select * from DUT"); err != nil {
// 		return dut, fmt.Errorf("failed to get dut list from DUT: %v", err)
// 	}
// 	return dut, nil
// }

// func (r *ResultMySQL) GetAllTest() ([]models.Test, error) {
// 	test := []models.Test{}
// 	if err := r.db.Select(test, "select * from Test"); err != nil {
// 		return test, fmt.Errorf("failed to get dut list from DUT: %v", err)
// 	}
// 	return test, nil
// }

// To-do: join Table?
func (r *ResultMySQL) SearchTestHaus(params utils.QueryParameter) ([]models.RawDataFromResult, error) {
	output := []models.RawDataFromResult{}
	options, err := utils.ToConditions(params)
	if err != nil {
		return output, fmt.Errorf("failed to parse the searching input: %v", err)
	}
	sql := "select Result.*, DUT.board as board, DUT.model as model, Test.name as test_name from Result, Dut, Test Where Resut.dut_id = Dut.dut_id AND Resut.test_id = Test.test_id "
	args := make([]interface{}, 0, len(options))
	for _, option := range options {
		sql += " AND " + option.Where
		args = append(args, option.Value)
	}
	if err := r.db.Select(output, sql, args...); err != nil {
		return output, fmt.Errorf("failed to get result: %v", err)
	}
	return output, nil
}
