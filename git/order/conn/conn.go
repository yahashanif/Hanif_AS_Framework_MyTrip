package database

import (
	"database/sql"
	"fmt"
	"time"

	"Hanif_Aulia_Sabri-MyTrip/git/order/parser"

	//load mysql driver
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
)

type DatabaseConnection interface {
	New(fn string) (*DbConnection, error)
	Open() (*sql.DB, error)
	Close()
	GetRows(rows *sql.Rows) (map[int]map[string]string, error)
	GetFirstRow() (string, error)
	Query(sqlStringName string, args ...interface{}) (*sql.Rows, error)
	Exec(sqlStringName string, args ...interface{}) (int64, error)
	Queryf(sqlStringName string, args ...interface{}) (*sql.Rows, error)
	Execf(sqlStringName string, args ...interface{}) (int64, error)
	InsertGetLastId(sqlStringName string, args ...interface{}) (int64, error)
}

//Config is global config for database connection
type DbConnection struct {
	Type     string            `yaml:"Type"`
	URL      string            `yaml:"URL"`
	Username string            `yaml:"Username"`
	Password string            `yaml:"Password"`
	Host     string            `yaml:"Host"`
	Schema   string            `yaml:"Schema"`
	SQL      map[string]string `yaml:"SQLCommand"`
	Db       *sql.DB
	Tx       *sql.Tx
}

//var c.Db *sql.DB

func New(fn string) (*DbConnection, error) {
	var c DbConnection
	if err := parser.LoadYAML(&fn, &c); err != nil {
		return nil, err
	}

	if c.URL == "" {
		//key := strings.Repeat(strings.ToUpper(c.Username), 2)
		//if password, err := crypt.TripleDesDecrypt(c.Password, []byte(key), crypt.PKCS5UnPadding); err != nil {
		//	if _, e := crypt.TripleDesEncrypt(c.Password, []byte(key), crypt.PKCS5Padding); e == nil {
		//		//log.Debugf("Decryption error, try %s instead\n", encpass)
		//	}
		//	return nil, err
		//} else {
		//	c.Password = password
		//}
	}

	return &c, nil
}

func (c DbConnection) Begin() error {
	var e error
	c.Tx, e = c.Db.Begin()
	if e != nil {
		return e
	}
	return nil
}

func (c DbConnection) Commit() error {
	return c.Tx.Commit()
}

func (c DbConnection) Rollback() error {
	return c.Tx.Rollback()
}

// OpenConnection prepares dbConnection for future connection to database
func (c DbConnection) Open() (*sql.DB, error) {
	if c.Username != "" && c.Password != "" && c.Host != "" && c.Schema != "" {
		c.URL = fmt.Sprintf("%s:%s@(%s)/%s", c.Username, c.Password, c.Host, c.Schema)
	}

	//fmt.Println(c.URL)

	c.Close()

	// Open database connection
	var err error
	//log.Debug("Initiating database connection ...")
	dbConn, err := sql.Open(c.Type, c.URL)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(100)

	dbConn.SetMaxIdleConns(20)
	dbConn.SetConnMaxLifetime(30 * time.Minute)

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// CloseConnection closes existing dbConnection
//
func (c DbConnection) Close() {
	if c.Db != nil {
		//log.Debug("Closing previous database connection.")
		c.Db.Close()
		c.Db = nil
	}
}

//ParsingRowsHelper parses recordset into map
func (c DbConnection) GetRows(rows *sql.Rows) (map[int]map[string]string, error) {
	var results map[int]map[string]string
	results = make(map[int]map[string]string)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	counter := 1
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// initialize the second layer
		results[counter] = make(map[string]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[counter][columns[i]] = value
		}
		counter++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//ParsingRowsAndGetValue parse and gets column value in first record
func (c DbConnection) GetFirstRow(rows *sql.Rows, key string) (string, error) {
	results, err := c.GetRows(rows)
	if err != nil {
		return "", err
	}
	return results[1][key], nil
}

// Query sends SELECT command to database
func (c DbConnection) Query(sqlStringName string, args ...interface{}) (*sql.Rows, error) {
	// if no dbConnection, return
	//
	if c.Db == nil {
		return nil, fmt.Errorf("Database needs to be initiated first.")
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	//fmt.Println(strSQL)

	rows, err := c.Db.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c DbConnection) QueryTx(sqlStringName string, args ...interface{}) (*sql.Rows, error) {
	// if no dbConnection, return
	//
	if c.Tx == nil {
		return nil, fmt.Errorf("Transaction needs to be initiated first.")
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	//fmt.Println(strSQL)

	rows, err := c.Tx.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//Exec executes UPDATE/INSERT/DELETE statements and returns rows affected
func (c DbConnection) Exec(sqlStringName string, args ...interface{}) (int64, error) {
	// if no dbConnection, return
	//
	if c.Db == nil {
		return 0, fmt.Errorf("Please OpenConnection prior Query")
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	// Execute the query
	res, err := c.Db.Exec(strSQL, args...)
	if err != nil {
		return 0, err //panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (c DbConnection) ExecTx(sqlStringName string, args ...interface{}) (int64, error) {
	// if no dbConnection, return
	//
	if c.Tx == nil {
		return 0, fmt.Errorf("Please Begin() transaction first")
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	// Execute the query
	res, err := c.Tx.Exec(strSQL, args...)
	if err != nil {
		return 0, err //panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (c DbConnection) InsertGetLastId(sqlStringName string, args ...interface{}) (int64, error) {
	// if no dbConnection, return
	//
	if c.Db == nil {
		return 0, fmt.Errorf("Please OpenConnection prior Query")
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	// Execute the query
	res, err := c.Db.Exec(strSQL, args...)
	if err != nil {
		//panic(err.Error()) // proper error handling instead of panic in your app
		return 0, err
	}

	rows, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (c DbConnection) Queryf(sql string, a ...interface{}) (*sql.Rows, error) {
	return c.Query(fmt.Sprintf(sql, a...))
}

func (c DbConnection) Execf(sql string, a ...interface{}) (int64, error) {
	return c.Exec(fmt.Sprintf(sql, a...))
}
