package db

import (
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/dchest/stemmer/porter2"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

func init() {
	cacheChannel = make(chan []osu.BeatmapSetsIN)
	go func() {
		for ins := range cacheChannel {
			insertStringIndex(ins)
		}
	}()

}

var (
	regexpReplace, _ = regexp.Compile(`[^0-9A-z]|[\[\]]`)
	cacheChannel     chan []osu.BeatmapSetsIN
)

func InsertCache(data []osu.BeatmapSetsIN) {
	cacheChannel <- data

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

func insertStringIndex(data []osu.BeatmapSetsIN) {
	defer func() {
		err, e := recover().(error)
		if e {
			pterm.Error.Println(err)
		}
	}()
	var insertDataa insertData
	//pterm.Info.Println(unsafe.Pointer(&insertDataa), string(*utils.ToJsonString(insertDataa)))

	for _, in := range data {

		artist := splitString(*in.Artist)
		creator := splitString(*in.Creator)
		title := splitString(*in.Title)
		tags := splitString(*in.Tags)
		insertDataa.Artist = append(insertDataa.Artist, row{
			KEY:          artist,
			BeatmapsetId: in.Id,
		})
		insertDataa.Creator = append(insertDataa.Creator, row{
			KEY:          creator,
			BeatmapsetId: in.Id,
		})
		insertDataa.Title = append(insertDataa.Title, row{
			KEY:          title,
			BeatmapsetId: in.Id,
		})
		insertDataa.Tags = append(insertDataa.Tags, row{
			KEY:          tags,
			BeatmapsetId: in.Id,
		})
		insertDataa.Strbuf = append(insertDataa.Strbuf, artist...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, creator...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, title...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, tags...)

	}
	pterm.Info.Println(string(*utils.ToJsonString(insertDataa))[:100])

	err := BulkInsertLimiter(
		"INSERT INTO SEARCH_CACHE_STRING_INDEX (STRING) VALUES %s ON DUPLICATE KEY UPDATE TMP = 1;",
		"(?)",
		utils.MakeArrayUnique(&insertDataa.Strbuf),
	)
	if err == nil {
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_ARTIST (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE TMP = 1;",
			"((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertDataa.Artist),
		)

		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_TITLE (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE TMP = 1;",
			"((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertDataa.Title),
		)
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_CREATOR (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE TMP = 1;",
			"((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertDataa.Creator),
		)
		_ = BulkInsertLimiter(
			"INSERT INTO SEARCH_CACHE_TAG (INDEX_KEY,BEATMAPSET_ID) VALUES %s ON DUPLICATE KEY UPDATE TMP = 1;",
			"((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)",
			toIndexKV(insertDataa.Tags),
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
