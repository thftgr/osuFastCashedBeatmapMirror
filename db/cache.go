package db

import (
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/dchest/stemmer/porter2"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

var (
	//STRING_INDEX     = map[string]*struct{}{}
	regexpReplace, _ = regexp.Compile(`[\[\]{}<>()]|^[,./\\;:]|[,./\\;:]$|\s`)
)

func LoadCache() {
	//rows, err := Maria.Query(`select STRING from SEARCH_CACHE_STRING_INDEX`)
	//if err != nil && err != sql.ErrNoRows {
	//	pterm.Error.Println(err)
	//	return
	//}
	//defer rows.Close()
	//var tmp string
	//for rows.Next() {
	//	err := rows.Scan(&tmp)
	//	if err != nil {
	//		pterm.Error.Println(err)
	//		continue
	//	}
	//	STRING_INDEX[tmp] = &struct{}{}
	//}
	//pterm.Info.Println("LoadCache() end")
}

func InsertCache(data *osu.BeatmapSetsIN) {
	//var ch = make(chan struct{})
	//go func() {
	var artistBuffer []string
	for _, s := range strings.Split(strings.ToLower(*data.Artist), " ") {
		s = strings.TrimSpace(regexpReplace.ReplaceAllString(s, " "))
		if s != "" {
			//if STRING_INDEX[s] == nil && s != "" {
			artistBuffer = append(artistBuffer, s)
			artistBuffer = append(artistBuffer, porter2.Stemmer.Stem(s))
		}
	}
	if len(artistBuffer) < 1 {
		return
	}
	keys := make(map[string]struct{})
	res := make([]interface{}, 0)
	for _, s := range artistBuffer {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	if len(res) < 1 {
		return
	}

	query := `INSERT IGNORE INTO osu.SEARCH_CACHE_STRING_INDEX (STRING) VALUES ` + strings.Join(repeatStringArray("(?)", len(res)), ",")
	_, err := Maria.Exec(query, res...)
	if err != nil {
		pterm.Error.Println(err, query, res)
	}

}
func repeatStringArray(s string, count int) (arr []string) {
	for i := 0; i < count; i++ {
		arr = append(arr, s)
	}
	return
}
