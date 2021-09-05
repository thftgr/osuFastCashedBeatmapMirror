package src

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
)

var Maria *sql.DB

func ConnectMaria(c *pterm.SpinnerPrinter) {

	db, err := sql.Open("mysql", Setting.Sql.Id+":"+Setting.Sql.Passwd+"@tcp("+Setting.Sql.Url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		c.UpdateText("RDBMS connected")

		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;");err != nil {
			c.Fail("SET SQL_SAFE_UPDATES FAIL.",err)
			panic(err)
		}
		c.Success("RDBMS Connected.")
	} else {
		c.Fail("RDBMS Connect Fail",err)
		panic(err)
	}
}



