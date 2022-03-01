package utils

import (
	"encoding/json"
	"log"
)

func ToJsonIndentString(i interface{}) (str *[]byte) {
	b, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	return &b
}
func ToJsonString(i interface{}) (str *[]byte) {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}
	return &b
}
