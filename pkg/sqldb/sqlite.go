package sqldb

import (
	"database/sql"
	"fmt"
)

type DB struct {
	pathToDB string
	dbName   string
	conn     *sql.DB
}

func NewDB(pathToDB, dbName string) *DB {
	return &DB{
		pathToDB: pathToDB,
		dbName:   dbName,
	}
}

func (d *DB) Query(query string) (*sql.Rows, error) {
	d.connect()
	defer d.close()
	res, err := d.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return res, nil
}

func (d *DB) connect() (bool, error) {
	var err error
	d.conn, err = sql.Open("sqlite3", d.pathToDB+"/"+d.dbName)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("db connect(open): %w", err)
	}

	if err = d.conn.Ping(); err != nil {
		return false, fmt.Errorf("db connect(ping): %w", err)
	}

	return true, nil
}

func (d *DB) close() error {
	err := d.conn.Close()
	if err != nil {
		return fmt.Errorf("db close: %w", err)
	}

	return nil
}
