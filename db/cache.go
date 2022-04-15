package db

import (
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/dchest/stemmer/porter2"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

var (
	regexpReplace, _ = regexp.Compile(`[^0-9A-z]|[\[\]]`)
)

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
	//return
	defer func() {
		err, e := recover().(error)
		if e {
			pterm.Error.Println(err)
		}
	}()
	insertData := insertData{}

	for _, in := range *data {
		artist := splitString(*in.Artist)
		creator := splitString(*in.Creator)
		title := splitString(*in.Title)
		tags := splitString(*in.Tags)
		insertData.Artist = append(insertData.Artist, row{
			KEY:          artist,
			BeatmapsetId: in.Id,
		})
		insertData.Creator = append(insertData.Creator, row{
			KEY:          creator,
			BeatmapsetId: in.Id,
		})
		insertData.Title = append(insertData.Title, row{
			KEY:          title,
			BeatmapsetId: in.Id,
		})
		insertData.Tags = append(insertData.Tags, row{
			KEY:          tags,
			BeatmapsetId: in.Id,
		})
		insertData.Strbuf = append(insertData.Strbuf, artist...)
		insertData.Strbuf = append(insertData.Strbuf, creator...)
		insertData.Strbuf = append(insertData.Strbuf, title...)
		insertData.Strbuf = append(insertData.Strbuf, tags...)

	}

	err := BulkInsertLimiter(
		"INSERT INTO SEARCH_CACHE_STRING_INDEX (STRING) VALUES %s ON DUPLICATE KEY UPDATE STRING= VALUES(STRING);",
		"(?)",
		makeArrayUnique(insertData.Strbuf),
	)
	if err == nil {
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_ARTIST (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE BEATMAPSET_ID= VALUES(BEATMAPSET_ID) ;",
			"((SELECT ID FROM  SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertData.Artist),
		)

		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_TITLE (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE BEATMAPSET_ID= VALUES(BEATMAPSET_ID) ;",
			"((SELECT ID FROM  SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertData.Title),
		)
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_CREATOR (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE BEATMAPSET_ID= VALUES(BEATMAPSET_ID) ;",
			"((SELECT ID FROM  SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertData.Creator),
		)
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_TAG (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE BEATMAPSET_ID= VALUES(BEATMAPSET_ID) ;",
			"((SELECT ID FROM  SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertData.Tags),
		)
	}

}

func toIndexKV(data []row) (AA []interface{}) {
	for _, A := range data {
		for _, K := range A.KEY {
			AA = append(AA, K, A.BeatmapsetId)
		}
	}
	return
}

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(regexpReplace.ReplaceAllString(input, " ")), " ") {
		if s == "" || s == " " {
			continue
		}
		ss = append(ss, s, porter2.Stemmer.Stem(s))
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
