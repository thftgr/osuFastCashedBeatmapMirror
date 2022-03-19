package main

import (
	"github.com/dchest/stemmer/porter2"
	"log"
	"regexp"
	"strings"
)

func main() {
	title := "Kotodama Refrain. (katagiri Bootleg) : ;"
	reg, _ := regexp.Compile(`[\[\]{}<>()]|^[,./\\;:]|[,./\\;:]$|\s`)

	for _, s := range strings.Split(title, " ") {
		s = strings.TrimSpace(reg.ReplaceAllString(s, " "))
		if s == "" {
			continue
		}
		log.Println(s)
		log.Println(porter2.Stemmer.Stem(s))

	}

}
