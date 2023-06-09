package result

import (
	"backend/lib/models"
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
