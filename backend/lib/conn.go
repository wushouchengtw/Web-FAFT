package lib

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(cfg *Database) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Scheme,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxConnection)
	return db, nil
}
