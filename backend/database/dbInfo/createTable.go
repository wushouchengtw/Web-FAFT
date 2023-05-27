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
	TableName TableName `json:"tablename"`
	Sql       []KeyType `json:"sql"`
}

type TableName string

type TableJson struct {
	JsonFile string
}

func (s *TableJson) GenerateTableMap() map[TableName]string {
	out, _ := ioutil.ReadFile(s.JsonFile)
	var t []TableConfig
	if err := json.Unmarshal(out, &t); err != nil {
		log.Fatalf("Failed to unmarshal: %s", out)
	}

	tableMap := map[TableName]string{}
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
