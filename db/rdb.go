package db

import (
	"database/sql"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Maria *sql.DB
var Gorm *gorm.DB

func ConnectRDBMS() {

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
		//                                        config.Config.Debug ? debug : error
		Logger:                                   logger.Default.LogMode(utils.TernaryOperator(config.Config.Debug, logger.Info, logger.Error)),
		CreateBatchSize:                          100,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if Gorm = orm; orm != nil {
		pterm.Success.Println("RDBMS orm connected")
	} else {
		pterm.Error.Println("RDBMS orm Connect Fail", err)
		panic(err)
	}

}
