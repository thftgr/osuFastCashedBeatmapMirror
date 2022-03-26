package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func main() {
	//title := "Kotodama Refrain. (katagiri Bootleg) : ;"
	//reg, _ := regexp.Compile(`[\[\]{}<>()]|^[,./\\;:]|[,./\\;:]$|\s`)
	//
	//for _, s := range strings.Split(title, " ") {
	//	s = strings.TrimSpace(reg.ReplaceAllString(s, " "))
	//	if s == "" {
	//		continue
	//	}
	//	log.Println(s)
	//	log.Println(porter2.Stemmer.Stem(s))
	//
	//}
	fmt.Println(bulkInsertLimiter(
		"insert into (A,B,C,D) VALUES %s ;",
		4,
		[]interface{}{
			"1", "2",
			"3", "4",
			"5", "6",
			"7", "8",
			"9", "10",
			"11", "12",
			"13", "14",
			"15", "16",
			"17", "18",
			"19", "20",
		},
	))

}

func bulkInsertLimiter(query string, valueSize int, data []interface{}) (err error) {
	dataSize := len(data)

	if dataSize%valueSize != 0 {
		return errors.New(fmt.Sprintf("dataSize \%valueSize != 0"))
	}
	values := "(" + strings.Join(stringRepeatArray("?", valueSize), ",") + ")"
	size := valueSize * 2

	var j int
	for i := 0; i < dataSize; i += size {
		j += size
		if j > dataSize {
			j = dataSize
		}

		fmt.Println(fmt.Sprintf(query, strings.Join(stringRepeatArray(values, len(data[i:j])/valueSize), ",")), data[i:j])
	}
	return

}
func stringRepeatArray(s string, count int) (arr []string) {
	for i := 0; i < count; i++ {
		arr = append(arr, s)
	}
	return
}
