package dbInfo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type KeyType struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type TableConfig struct {
	TableName tableName `json:"tablename"`
	Sql       []KeyType `json:"sql"`
}

type tableName string

type tableJson struct {
	jsonFile string
}

func (s *tableJson) GenerateTableMap() map[tableName]string {
	out, _ := ioutil.ReadFile(s.jsonFile)
	var t []TableConfig
	if err := json.Unmarshal(out, &t); err != nil {
		log.Fatalf("Failed to unmarshal: %s", out)
	}

	tableMap := map[tableName]string{}
	for _, table := range t {
		var sqlStmt string
		for index, keyMap := range table.Sql {
			if index == len(table.Sql)-1 {
				sqlStmt += keyMap.Key + " " + keyMap.Type
			} else {
				sqlStmt += keyMap.Key + " " + keyMap.Type + ","
			}
		}
		tableMap[table.TableName] = sqlStmt
	}
	return tableMap
}

func CreateTable(DB *sql.DB, tableName, keyMap string) error {
	sql := `CREATE TABLE IF NOT EXISTS %s(`
	sql += keyMap + `); `

	sqlStmt := fmt.Sprintf(sql, tableName)
	return RunSqlStmt(DB, sqlStmt)
}

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
