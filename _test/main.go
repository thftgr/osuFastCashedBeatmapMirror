package main

import (
	"fmt"
	"github.com/dchest/stemmer/porter2"
	"regexp"
	"strings"
)

var (
	//STRING_INDEX     = map[string]*struct{}{}
	regexpReplace, _ = regexp.Compile(`[^0-9A-z]|[\[\]]`)
)

func main() {

	//fmt.Println(BulkInsertLimiter(
	//	"insert into (A,B,C,D) VALUES %s ;",
	//	"(?)",
	//	[]interface{}{
	//		"1", "2",
	//		"3", "4",
	//		"5", "6",
	//		"7", "8",
	//		"9", "10",
	//		"11", "12",
	//		"13", "14",
	//		"15", "16",
	//		"17", "18",
	//		"19", "20",
	//	}...,
	//))
	fmt.Println(splitString("[Asakura] (Kirika) (CV: Akabane Kyouko)@"))

}

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(regexpReplace.ReplaceAllString(input, " ")), " ") {

		if s == "" {
			continue
		} else if strings.Contains(s, " ") {
			ss = append(ss, splitString(s)...)
		} else {
			ss = append(ss, s)
			ss = append(ss, porter2.Stemmer.Stem(s))
		}
	}
	return
}
