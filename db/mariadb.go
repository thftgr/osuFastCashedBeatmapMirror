package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/config"
)

var Maria *sql.DB

func ConnectMaria() {

	db, err := sql.Open("mysql", config.Setting.Sql.Id+":"+config.Setting.Sql.Passwd+"@tcp("+config.Setting.Sql.Url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		pterm.Success.Println("RDBMS connected")

		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;"); err != nil {
			pterm.Error.Println("SET SQL_SAFE_UPDATES FAIL.")
			panic(err)
		}
		//pterm.Success.Println("RDBMS Connected.")
	} else {
		pterm.Error.Println("RDBMS Connect Fail", err)
		panic(err)
	}
}
