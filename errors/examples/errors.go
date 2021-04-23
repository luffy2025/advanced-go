package examples

import (
	"database/sql"
	"github.com/pkg/errors"
)

func errNoRows() ([]int, error) {
	db, err := sql.Open("mysql", "user:passwd@tcp(localhost:3306)/testdb")
	if err != nil {
		return nil, errors.Wrapf(err, "open mysql db:%s error", "testdb")
	}
	defer db.Close()

	sqlStr := "SELECT `id` from `students` WHERE `age` < ?"
	rows, err := db.Query(sqlStr, 10)
	if err != nil {
		if err == sql.ErrNoRows {
			return []int{}, nil
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
