package dbInfo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type SqlDB struct {
	// userName for connecting sql DB
	UserName string
	// password for connecting sql DB
	Password string
	// IP for the sql server
	Ip string
	// db name used for the connection
	DbName string
}

type weekDay struct {
	startDay string
	endDay   string
}

// Table name in SQL
const (
	StainlessTable = "Stainless_Result"
	CaseTable      = "CaseTable"
	LocalTestTable = "Result"
)

func NewSql(dbConfig DbConfig) *SqlDB {
	return &SqlDB{
		UserName: dbConfig.UserName,
		Password: dbConfig.Password,
		Ip:       dbConfig.Ip,
		DbName:   dbConfig.Name,
	}
}

func (db *SqlDB) NewSqlDBConnect() (*sql.DB, error) {
	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", db.UserName, db.Password, db.Ip, db.DbName))
	if err != nil {
		return nil, errors.New("failed to open sql DB")
	}
	return DB, nil
}

func RunSqlStmt(DB *sql.DB, sql string, value ...interface{}) error {
	sqlStmt, err := DB.Prepare(sql)
	if err != nil {
		log.Fatal("Failed to prepare sql: ", err)
	}
	defer sqlStmt.Close()

	if _, err := sqlStmt.Exec(value...); err != nil {
		log.Fatal("Failed to run sql: ", err)
	}
	return nil
}
