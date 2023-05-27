package dbInfo

import (
	"database/sql"
	"fmt"
)

func CreateTableCaseTable(DB *sql.DB) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id int AUTO_INCREMENT PRIMARY KEY,
		CaseName varchar(255),
		Owner varchar(30),
		Ticket varchar(20)
	); `, CaseTable)
	return RunSqlStmt(DB, sql)
}

func CreateTableLocalTestTable(DB *sql.DB) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id int AUTO_INCREMENT PRIMARY KEY,
		Time datetime,
		tester varchar(15),
		name varchar(30),
		board varchar(20),
		model varchar(20),
		version varchar(20),
		logPath varchar(40),
		result varchar(5),
		reason text
	); `, LocalTestTable)
	return RunSqlStmt(DB, sql)
}

func CreateTableStainless(DB *sql.DB) error {
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id int AUTO_INCREMENT PRIMARY KEY,
			time TIMESTAMP,
			duration varchar(10),
			suite varchar(40),
			board varchar(20),
			model varchar(20), 
			buildVersion varchar(20),
			host varchar(60),
			testName varchar(80),
			status varchar(10),
			reason blob,
			firmwareROVersion varchar(50),
			firmwareRWVersion varchar(50)
        ); `, StainlessTable)
	return RunSqlStmt(DB, sql)
}
