package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"regexp"
	"time"
)

type InsertQueue struct {
	Ch    chan error
	Query string
	Args  []any
}

var InsertQueueChannel chan InsertQueue
var queryNameRegex, _ = regexp.Compile("^(/[*])(.+?)([*]/)")

func AddInsertQueue(query string, args ...any) {
	InsertQueueChannel <- InsertQueue{
		Query: query,
		Args:  args,
	}
}

func init() {
	go func() {
		for {
			InsertQueueChannel = make(chan InsertQueue)
			for ins := range InsertQueueChannel {
				st := time.Now().UnixMilli()
				result := Gorm.Exec(ins.Query, ins.Args...)
				err := result.Error
				ins := ins

				go func() { //로그용
					queryName := fmt.Sprintf("%s | %-50s | ", time.Now().Format("15:04:05.000"), pterm.Yellow(queryNameRegex.FindString(ins.Query[:100])))

					if err == nil && result.RowsAffected > 0 {
						pterm.Info.Printfln(queryName+"%5d ROWS | %6dms |", result.RowsAffected, time.Now().UnixMilli()-st)
					}
					if ins.Ch != nil {
						if err != nil {
							ins.Ch <- errors.New(queryName + err.Error() + ins.Query)
						} else {
							ins.Ch <- nil
						}
					} else if err != nil {
						pterm.Error.Println(queryName+err.Error(), ins.Query)
					}
				}()
			}
		}
	}()
}
