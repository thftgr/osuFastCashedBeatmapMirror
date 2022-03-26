package db

import (
	"errors"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/dchest/stemmer/porter2"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

var (
	//STRING_INDEX     = map[string]*struct{}{}
	regexpReplace, _ = regexp.Compile(`^[,./\\;:\[\]{}<>()]|[,./\\;:\[\]{}<>()]$|\s`)
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

func InsertCache(data *[]osu.BeatmapSetsIN) {
	go insertStringIndex(data)

}

type insertData struct {
	Strbuf  []string
	Artist  []row
	Creator []row
	Title   []row
	Tags    []row
}
type row struct {
	KEY          []string
	BeatmapsetId int
}

func insertStringIndex(data *[]osu.BeatmapSetsIN) {
	defer func() {
		err, e := recover().(error)
		if e {
			pterm.Error.Println(err)
		}
	}()
	insertData := insertData{}

	for _, in := range *data {

		artist = append(artist, splitString(*in.Artist)...)
		creator = append(creator, splitString(*in.Creator)...)
		title = append(title, splitString(*in.Title)...)
		tags = append(tags, splitString(*in.Tags)...)

	}
	stringBuf = append(stringBuf, artist...)
	stringBuf = append(stringBuf, creator...)
	stringBuf = append(stringBuf, title...)
	stringBuf = append(stringBuf, tags...)

	if len(stringBuf) < 1 {
		return
	}
	res := makeArrayUnique(stringBuf)
	if len(res) < 1 {
		return
	}
	query := `INSERT IGNORE INTO osu.SEARCH_CACHE_STRING_INDEX (STRING) VALUES ` + strings.Join(repeatStringArray("(?)", len(res)), ",") + ";"
	_, err := Maria.Exec(query, res...)
	if err != nil {
		pterm.Error.Println(err, query, res)
		return
	}

}

func bulkInsertLimiter(query, values string, data []interface{}) (err error) {
	dataSize := len(data)
	varSize := strings.Count(values, "?")
	if dataSize < 1 {
		return
	}
	if dataSize%varSize != 0 {
		return errors.New("args length not match")
	}
	strings.Repeat()

}

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(input), " ") {
		s = strings.TrimSpace(regexpReplace.ReplaceAllString(s, " "))
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

func repeatStringArray(s string, count int) (arr []string) {
	for i := 0; i < count; i++ {
		arr = append(arr, s)
	}
	return
}
func makeArrayUnique(array []string) []interface{} {

	keys := make(map[string]struct{})
	res := make([]interface{}, 0)
	for _, s := range array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	return res
}
