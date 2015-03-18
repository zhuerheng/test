package main

import (
	"database/sql"
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
)

func main() {

	http.HandleFunc("/", do)
	http.ListenAndServe(":33333", nil)
}

func do(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:qazxs913@tcp(adserver1:3306)/adserver-staging")
	db.SetMaxOpenConns(2000000)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * from qtad_banners where description = ?", string(b))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	column, _ := rows.Columns()
	scanArgs := make([]interface{}, len(column))
	values := make([]sql.NullString, len(column))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	var s string

	for rows.Next() {
		rows.Scan(scanArgs...)
		for i, col := range values {
			if col.Valid {
				s = s + column[i] + ": " + col.String + "\n"
			} else {
				s = s + column[i] + ": \n"
			}
		}
		s = s + "----------------------------\n"
	}
	http.Error(w, s, 200)
}
