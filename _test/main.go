package main

import (
	"database/sql"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/pterm/pterm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	ii := []int{0, 1, 2, 3, 4, 5, 5, 5}
	fmt.Println(ii)

	ii = utils.MakeArrayUnique(ii)
	fmt.Println(ii)

	//ConnectMaria()
	//str := Gorm.Select("ID").Table("SEARCH_CACHE_STRING_INDEX").Where("STRING IN (?)", []string{"my", "love"})
	//C1 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_TITLE").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//C2 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_ARTIST").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//C3 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_CREATOR").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//C4 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_TAG").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//
	//Gorm.Select("BEATMAPSET_ID").Table("(? UNION ALL ? UNION ALL ? UNION ALL ?) A;", C1, C2, C3, C4)
	//
	//Gorm.ToSQL(C1)

	//	//var i int
	//	////SELECT BEATMAPSET_ID from (
	//	////      SELECT BEATMAPSET_ID from SEARCH_CACHE_TITLE
	//	////      where (0xff & 1 = 1) AND INDEX_KEY in ( select ID from SEARCH_CACHE_STRING_INDEX where STRING in( 'my','love' ))
	//	////      GROUP BY BEATMAPSET_ID having count(*) >= 2
	//	////      UNION ALL SELECT BEATMAPSET_ID from SEARCH_CACHE_ARTIST
	//	////      where (0xff & 2 = 2) AND INDEX_KEY in ( select ID from SEARCH_CACHE_STRING_INDEX where STRING in( 'my','love' ))
	//	////      GROUP BY BEATMAPSET_ID having count(*) >= 2
	//	////      UNION ALL SELECT BEATMAPSET_ID from SEARCH_CACHE_CREATOR
	//	////      where (0xff & 4 = 4) AND INDEX_KEY in ( select ID from SEARCH_CACHE_STRING_INDEX where STRING in( 'my','love' ))
	//	////      GROUP BY BEATMAPSET_ID having count(*) >= 2
	//	////      UNION ALL SELECT BEATMAPSET_ID from SEARCH_CACHE_TAG
	//	////      where (0xff & 8 = 8) AND INDEX_KEY in ( select ID from SEARCH_CACHE_STRING_INDEX where STRING in( 'my','love' ))
	//	////      GROUP BY BEATMAPSET_ID having count(*) >= 2
	//	////) A
	//	////;
	//	//st := time.Now().UnixMicro()
	//	//et := time.Now().UnixMicro()
	//	//time.Sleep(time.Second * 5)
	//	//st = time.Now().UnixMicro()
	//	//_ = Maria.QueryRow(`SELECT BEATMAPSET_ID FROM (
	//	//          SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TITLE WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_ARTIST WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_CREATOR WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TAG WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//) A; `).Scan(&i)
	//	//et = time.Now().UnixMicro()
	//	//fmt.Println(et-st, i)
	//	//
	//	//st = time.Now().UnixMicro()
	//	//Gorm.Raw(`SELECT BEATMAPSET_ID FROM (
	//	//          SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TITLE WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_ARTIST WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_CREATOR WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TAG WHERE INDEX_KEY IN (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN ('my','love')) GROUP BY BEATMAPSET_ID HAVING count(*) >= 2
	//	//) A; `).Scan(&i)
	//	//et = time.Now().UnixMicro()
	//	//fmt.Println(et-st, i)
	//	//
	//	//st = time.Now().UnixMicro()
	//	//str := Gorm.Select("ID").Table("SEARCH_CACHE_STRING_INDEX").Where("STRING IN (?)", []string{"my", "love"})
	//	//C1 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_TITLE").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//	//C2 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_ARTIST").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//	//C3 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_CREATOR").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//	//C4 := Gorm.Select("BEATMAPSET_ID").Table("SEARCH_CACHE_TAG").Where("INDEX_KEY IN (?)", str).Group("BEATMAPSET_ID").Having("count(*) >= ?", 2)
	//	//
	//	//Gorm.Select("BEATMAPSET_ID").Table("(? UNION ALL ? UNION ALL ? UNION ALL ?) A;", C1, C2, C3, C4).Scan(&i)
	//	//
	//	//et = time.Now().UnixMicro()
	//	//fmt.Println(et-st, i)
	//
}

var Maria *sql.DB
var Gorm *gorm.DB

func ConnectMaria() {

	db, err := sql.Open("mysql", "root:myaimgod!1?@tcp(192.168.0.50:3306)/osu")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)

		pterm.Success.Println("RDBMS connected")
		var i int
		Maria.QueryRow("SELECT 1;").Scan(&i)

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
		var i int
		Gorm.Raw("SELECT 1;").Scan(&i)
		pterm.Success.Println("RDBMS orm connected")
	} else {
		pterm.Error.Println("RDBMS orm Connect Fail", err)
		panic(err)
	}

}
