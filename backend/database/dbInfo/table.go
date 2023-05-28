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

func (s *TableJson) GenerateTableMap() (mapToCreateTables map[TableName]string, mapTableKeys map[TableName][]string) {
	out, _ := ioutil.ReadFile(s.JsonFile)
	var t []TableConfig
	if err := json.Unmarshal(out, &t); err != nil {
		log.Fatalf("Failed to unmarshal: %s", out)
	}

	tableCreateSql := map[TableName]string{}
	tableKeyContains := map[TableName][]string{}
	for _, table := range t {
		var sql string
		var keySlice []string
		for index, keyMap := range table.Sql {
			keySlice = append(keySlice, keyMap.Key)
			if index == len(table.Sql)-1 {
				sql += keyMap.Key + " " + keyMap.Type
			} else {
				sql += keyMap.Key + " " + keyMap.Type + ","
			}
		}
		sql = `CREATE TABLE IF NOT EXISTS %s(` + sql
		sql += `); `
		sqlStmt := fmt.Sprintf(sql, table.TableName)
		tableCreateSql[table.TableName] = sqlStmt
		tableKeyContains[table.TableName] = keySlice
	}
	return tableCreateSql, tableKeyContains
}

func CreateTable(DB *sql.DB, sqlStmt string) error {
	return RunSqlStmt(DB, sqlStmt)
}
