package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println("failed to open database:", err.Error())
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("%s error ping database: %s", err.Error())
		return
	}
	defer db.Close()

	name, errs := GetUserName(db, 11)

	if errs != nil {

		fmt.Println(errs)

		return
	}

	fmt.Println(name)

}

func GetUserName(db *sql.DB, id int) (string, error) {
	var name string
	sqlText := "SELECT `name` FROM dtk_test  where id = ?"
	err2 := db.QueryRow(sqlText, id).Scan(&name)

	if err2 != nil {

		//像上抛出错误,并且多的带上其他错误信息
		return name, errors.Wrap(err2, "查询出错 sql "+sqlText+" ")

	}
	return name, nil

}
