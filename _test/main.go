package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var queryNameRegex, _ = regexp.Compile("^(/[*])(.+?)([*]/)")

func main() {

	fmt.Println(strings.Replace("A               B", "  ", " ", -1))

	//log.Println(time.Now().Format("15:04:05.000"), queryNameRegex.FindString("/* INSERT SEARCH_CACHE_ARTIST */ INSERT INTO SEARCH_CACHE_ARTIST"))
}

func ToJsonString(i interface{}) (str string) {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}
	return string(b)
}
