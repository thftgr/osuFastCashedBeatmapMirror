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

func StringRepeatArray(s string, count int) (arr []string) {
	for i := 0; i < count; i++ {
		arr = append(arr, s)
	}
	return
}

func TernaryOperator(tf bool, t, f any) any {
	if tf {
		return t
	}
	return f
}
