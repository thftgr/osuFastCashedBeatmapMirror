package db

import (
	"errors"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/pterm/pterm"
	"strings"
)

//BulkInsertLimiter query = "INSERT INTO DB.TABLE (A,B,C,D)  VALUES %S ;"
func BulkInsertLimiter(query, value string, aa []interface{}) (err error) {
	dataSize := len(aa)
	valueSize := strings.Count(value, "?")
	if dataSize%valueSize != 0 {
		return errors.New(fmt.Sprintf("dataSize %% valueSize != 0"))
	}
	size := valueSize * 200

	var j int
	for i := 0; i < dataSize; i += size {
		j += size
		if j > dataSize {
			j = dataSize
		}
		if _, err := Maria.Exec(fmt.Sprintf(query, strings.Join(utils.StringRepeatArray(value, len(aa[i:j])/valueSize), ",")), aa[i:j]...); err != nil {
			pterm.Error.Println("eaa", err)
		}
	}
	return

}
