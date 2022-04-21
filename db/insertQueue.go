package db

import (
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
)

type InsertQueue struct {
	Ch    chan error
	Query string
	Args  []any
}

var InsertQueueChannel chan InsertQueue

func init() {
	InsertQueueChannel = make(chan InsertQueue)
	go func() {

		for ins := range InsertQueueChannel {
			result := Gorm.Exec(ins.Query, ins.Args...)
			err := result.Error
			if err == nil && result.RowsAffected > 0 {
				pterm.Info.Printfln("UPDATED %d ROWS", result.RowsAffected)
			}
			if ins.Ch != nil {
				if err != nil {
					ins.Ch <- errors.New(err.Error() + ins.Query)
				} else {
					ins.Ch <- nil
				}
			} else if err != nil {
				pterm.Error.Println(err, ins.Query)
			}
		}
	}()
}
