package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"time"
)

var Maria *sql.DB



func QueryOnly(sql string, parm ...interface{}) error {

	raws, err := Maria.Query(sql, parm...)
	if err != nil {
		if raws != nil {
			_ = raws.Close()
		}
		return err
	}

	return raws.Close()
}

func ConnectMaria() {

	db, err := sql.Open("mysql", src.Setting.Sql.Id+":"+src.Setting.Sql.Passwd+"@tcp("+src.Setting.Sql.Url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		fmt.Println("mariaDB connected")

		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;");err != nil {
			fmt.Println("SET SQL_SAFE_UPDATES FAIL")
			panic(err)
		}
		fmt.Println("SET SQL_SAFE_UPDATES PASS")
	} else {
		panic(err)
	}
}

func Upsert(query string, data []interface{}) {
	data = append(data,data[1:]...)
	err := QueryOnly(
		query,
		data...
	)
	if err != nil {
		fmt.Println(err,query)
		fmt.Println(err,fmt.Sprintf("%v",data))
	}
}

func ToDateTime(t interface{}) string {
	if t == nil {
		return "0000-00-00T00:00:00"
	}
	myDate, _ := time.Parse("2006-01-02T15:04:05-07:00", t.(string))
	return myDate.Format("2006-01-02T15:04:05")
}

func InsertAPILog(s ...interface{}) (err error) {
	rows, err := Maria.Query(QueryAPILog, s...)
	if err != nil {
		return
	}
	defer rows.Close()
	return
}


