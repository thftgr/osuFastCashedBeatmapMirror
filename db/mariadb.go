package db

import (
	"database/sql"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Maria *sql.DB
var Gorm *gorm.DB

//TODO my.cnf 에 innodb_autoinc_lock_mode=0 추가해야함
//TODO
//TODO
//TODO
func ConnectMaria() {

	db, err := sql.Open("mysql", config.Config.Sql.Url)
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)

		pterm.Success.Println("RDBMS connected")

		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;"); err != nil {

			pterm.Error.Println("SET SQL_SAFE_UPDATES FAIL.")
			panic(err)
		}
	} else {
		pterm.Error.Println("RDBMS Connect Fail", err)
		panic(err)
	}

	orm, err := gorm.Open(mysql.New(mysql.Config{Conn: Maria}), &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            logger.Default.LogMode(logger.Info),
		CreateBatchSize:   100,
	})
	if Gorm = orm; orm != nil {
		pterm.Success.Println("RDBMS orm connected")
	} else {
		pterm.Error.Println("RDBMS orm Connect Fail", err)
		panic(err)
	}

}
