package shallalist

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Tma123321"
	dbName := "tma"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(localhost:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func selectHosts(db *sql.DB, ini, fin int) (*[]string, error) {
	hosts := make([]string, fin)

	err := db.Ping()
	if err != nil {
		return nil, err
	}
	res, err := db.Query("SELECT DISTINCT(host) FROM tma_data WHERE NOT IS_IPV4(host) AND HOST != '' LIMIT ? OFFSET ?", fin, ini)
	defer res.Close()

	if err != nil {
		return nil, err
	}
	i := 0
	for res.Next() {
		var host string
		err := res.Scan(&host)
		if err != nil {
			return nil, err
		}
		if len(host) > 1 {
			if host[len(host)-1] == '/' {
				host = host[0 : len(host)-1]
			}
			hosts[i] = host
			i++
		}
	}
	return &hosts, nil
}

func insertCategory(db *sql.DB, host, category string, matches int) error {
	//here we will pass the db client as parameter
	//fmt.Printf("%v %v %v\n", host, category, matches)
	insForm, err := db.Prepare("INSERT INTO h_categories (host, category, matches) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	//fmt.Printf("INSERT INTO h_categories (host, category, matches) VALUES(%v,%v,%v", host, category, matches)
	_, err = insForm.Exec(host, category, matches)
	if err != nil {
		return err
	}
	return nil
}
