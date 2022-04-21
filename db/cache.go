package db

import (
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/dchest/stemmer/porter2"
	"github.com/pterm/pterm"
	"regexp"
	"strconv"
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
	Other   []row
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
			KEY:          utils.MakeArrayUnique(&artist),
			BeatmapsetId: in.Id,
		})
		insertDataa.Creator = append(insertDataa.Creator, row{
			KEY:          utils.MakeArrayUnique(&creator),
			BeatmapsetId: in.Id,
		})
		insertDataa.Title = append(insertDataa.Title, row{
			KEY:          utils.MakeArrayUnique(&title),
			BeatmapsetId: in.Id,
		})
		insertDataa.Tags = append(insertDataa.Tags, row{
			KEY:          utils.MakeArrayUnique(&tags),
			BeatmapsetId: in.Id,
		})
		insertDataa.Strbuf = append(insertDataa.Strbuf, artist...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, creator...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, title...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, tags...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, strconv.Itoa(in.Id))
		for _, beatmapIN := range *in.Beatmaps {
			other := splitString(*beatmapIN.Version)
			insertDataa.Strbuf = append(insertDataa.Strbuf, other...)
			insertDataa.Strbuf = append(insertDataa.Strbuf, strconv.Itoa(beatmapIN.Id))
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          []string{strconv.Itoa(beatmapIN.Id)},
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          []string{strconv.Itoa(beatmapIN.BeatmapsetId)},
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          other,
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
		}

	}
	if idata := utils.MakeArrayUniqueInterface(&insertDataa.Strbuf); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_STRING_INDEX */ INSERT INTO SEARCH_CACHE_STRING_INDEX (STRING) VALUES "+utils.StringRepeatJoin("(?)", ",", len(idata))+" ON DUPLICATE KEY UPDATE TMP = 1;",
			idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Artist); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_ARTIST */ INSERT INTO SEARCH_CACHE_ARTIST (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Title); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_TITLE */ INSERT INTO SEARCH_CACHE_TITLE (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Creator); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_CREATOR */ INSERT INTO SEARCH_CACHE_CREATOR (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}
	if idata := toIndexKV(insertDataa.Tags); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_TAG */ INSERT INTO SEARCH_CACHE_TAG (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}
	if idata := toIndexKV(insertDataa.Other); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_OTHER */ INSERT INTO SEARCH_CACHE_OTHER (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
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
