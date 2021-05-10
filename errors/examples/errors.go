package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

// ErrDBNotFound 自定义未找到错误
var ErrDBNotFound = errors.New("db not found.")

func main() {
	//db, err := newMyDB("user:passwd@tcp(localhost:3306)/testdb")
	db, err := newMockMyDB()
	if err != nil {
		log.Fatalf("main: new mydb error. error:%+v", err)
	}
	defer db.Close()

	results, err := db.Query("SELECT `id` from `students` WHERE `age` < ?", 10)
	if err != nil {
		if errors.Is(err, ErrDBNotFound) {
			log.Println("db.Query(): results are empty.")
		} else {
			log.Printf("main: db error. error:%+v", err)
			return
		}
	}
	for _, r := range results {
		fmt.Println(r)
	}
}

type db interface {
	Close() error
	Query(sql string, args ...interface{}) ([]int, error)
}

type myDB struct {
	db *sql.DB
}

func newMyDB(dbsn string) (db, error) {
	db, err := sql.Open("mysql", dbsn)
	if err != nil {
		return nil, errors.Wrapf(err, "open mysql db:%s error", "testdb")
	}
	return &myDB{db: db}, nil
}

func (m *myDB) Close() error {
	return m.db.Close()
}

func (m *myDB) Query(sqlStr string, args ...interface{}) ([]int, error) {

	rows, err := m.db.Query(sqlStr, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []int{}, ErrDBNotFound
		}
		return nil, errors.Wrapf(err, "query mysql error. sql:%s", sqlStr)
	}

	var ids []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrapf(err, "scan rows error. sql:%s", sqlStr)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

type mockMyDB struct{}

func newMockMyDB() (db, error) {
	return &mockMyDB{}, nil
}

func (*mockMyDB) Close() error {
	return nil
}

func (*mockMyDB) Query(sqlStr string, args ...interface{}) ([]int, error) {
	return []int{}, ErrDBNotFound
}
