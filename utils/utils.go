package utils

import (
	"encoding/json"
	"log"
	"strings"
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

func TernaryOperator[V any](tf bool, T, F V) V {
	if tf {
		return T
	}
	return F
}

func MakeArrayUnique[T comparable](array *[]T) (res []T) {

	keys := make(map[T]struct{})
	for _, s := range *array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	return
}
func MakeArrayUniqueInterface[T comparable](array *[]T) (res []interface{}) {

	keys := make(map[T]struct{})
	for _, s := range *array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	return
}
func StringRepeatJoin(str, sep string, count int) string {
	var repeatBuf []string
	for i := 0; i < count; i++ {
		repeatBuf = append(repeatBuf, str)
	}
	return strings.Join(repeatBuf, sep)

}
