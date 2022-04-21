package Route

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (s *SearchQuery) parsePage() {
	// 에러 발생시 int value = 0
	page, _ := strconv.Atoi(s.Page)
	pageSize, _ := strconv.Atoi(s.PageSize)
	if page < 0 {
		page = 0
	}
	if pageSize < 10 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	if page == 0 && pageSize == 0 {
		s.Page = "LIMIT 50"
	} else if page != 0 && pageSize == 0 {
		s.Page = fmt.Sprintf("LIMIT %d,50", page*50)
	} else if page == 0 && pageSize != 0 {
		s.Page = fmt.Sprintf("LIMIT 0,%d", pageSize)
	} else if page != 0 && pageSize != 0 {
		s.Page = fmt.Sprintf("LIMIT %d,%d", page*pageSize, pageSize)
	}

	return
}

//

var (
	regexpReplace, _ = regexp.Compile(`[^0-9A-z]|[\[\]]`)
	mode             = map[string]int{
		"0": 0, "o": 0, "std": 0, "osu": 0, "osu!": 0, "standard": 0,
		"1": 1, "t": 1, "taiko": 1, "osu!taiko": 1,
		"2": 2, "c": 2, "ctb": 2, "catch": 2, "osu!catch": 2,
		"3": 3, "m": 3, "mania": 3, "osu!mania": 3,
	}
	ranked = map[string][]int{
		"ranked":    {1, 2},
		"qualified": {3},
		"loved":     {4},
		"pending":   {0},
		"wip":       {-1},
		"graveyard": {-2},
		"unranked":  {0, -1, -2},
		"-2":        {-2},
		"-1":        {-1},
		"0":         {0},
		"1":         {1},
		"2":         {2},
		"3":         {3},
		"4":         {4},
		"default":   {4, 2, 1},
	}
	orderBy = map[string]string{
		"ranked_asc":           "ranked_date",
		"ranked_date":          "ranked_date",
		"ranked_date asc":      "ranked_date",
		"favourites_asc":       "favourite_count",
		"favourite_count":      "favourite_count",
		"favourite_count asc":  "favourite_count",
		"plays_asc":            "play_count",
		"play_count":           "play_count",
		"play_count asc":       "play_count",
		"updated_asc":          "last_updated",
		"last_updated":         "last_updated",
		"last_updated asc":     "last_updated",
		"title_asc":            "title",
		"title":                "title",
		"title asc":            "title",
		"artist_asc":           "artist",
		"artist":               "artist",
		"artist asc":           "artist",
		"ranked_desc":          "ranked_date desc",
		"ranked_date desc":     "ranked_date desc",
		"favourites_desc":      "favourite_count desc",
		"favourite_count desc": "favourite_count desc",
		"plays_desc":           "play_count desc",
		"play_count desc":      "play_count desc",
		"updated_desc":         "last_updated desc",
		"last_updated desc":    "last_updated desc",
		"title_desc":           "title desc",
		"title desc":           "title desc",
		"artist_desc":          "artist desc",
		"artist desc":          "artist desc",
		"default":              "ranked_date desc",
	}
)

func (s *SearchQuery) parseNsfw() {
	ss := strings.ToLower(s.Nsfw)
	switch ss {
	case "1", "all":
		s.Nsfw = "all"
	default:
		s.Nsfw = "0"
	}
	return
}
func (s *SearchQuery) parseOption() {
	ss := strings.ToLower(s.Option)
	if ss == "" {
		s.OptionB |= 0xFF
		return
	}
	sss := strings.Split(ss, ",")
	for i := 0; i < len(sss); i++ {
		switch ss {
		case "artist", "a":
			s.OptionB |= 0x01
		case "creator", "c":
			s.OptionB |= 0x02
		case "tag", "tg":
			s.OptionB |= 0x04
		case "title", "t":
			s.OptionB |= 0x08
		}
	}
	return
}

func (s *SearchQuery) parseExtra() {
	s.Extra = strings.ToLower(strings.TrimSpace(s.Extra))
	if s.Storyboard != "1" && s.Storyboard != "all" {
		if strings.Contains(s.Extra, "storyboard") {
			s.Storyboard = "1"
		} else {
			s.Storyboard = "all"
		}
	}
	if s.Video != "1" && s.Video != "all" {
		if strings.Contains(s.Extra, "video") {
			s.Video = "1"
		} else {
			s.Video = "all"
		}
	}
}

type minMax struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

func (v *minMax) minAsString() string {
	return fmt.Sprintf("%.1f", v.Min)
}
func (v *minMax) maxAsString() string {
	return fmt.Sprintf("%.1f", v.Max)
}
func (v *minMax) minMaxAsQuery() (query string) {
	if v == nil || (v.Min == 0 && v.Max == 0) {
		return
	} else if v.Min != 0 && v.Max == 0 {
		return fmt.Sprintf(">= %.1f", v.Min)
	} else if v.Min == 0 && v.Max != 0 {
		return fmt.Sprintf("<= %.1f", v.Max)
	}
	return fmt.Sprintf("BETWEEN %.1f AND %.1f", v.Min, v.Max)
}
func (v *minMax) minMaxIsNil() (isNil bool) {
	return v == nil || (v.Min == 0 && v.Max == 0)
}
func (s *SearchQuery) parseQuery() {
	s.parseNsfw()
	s.parsePage()
	s.parseExtra()
	s.parseOption()

}
func (s *SearchQuery) parseRankedStatus() (status []int) {
	statuss := strings.Split(s.Ranked, ",")
	for _, st := range statuss {
		st = strings.ToLower(strings.TrimSpace(st))
		if st == "" {
			continue
		}
		rs := ranked[st]
		if len(rs) < 1 || rs == nil {
			return ranked["default"]
		}

	}
	statuss = utils.MakeArrayUnique(&statuss)
	return
}

type SearchQuery struct {
	// global
	Extra string `query:"e" json:"extra"` // 스토리보드 비디오.

	// set
	Ranked     string `query:"s" json:"ranked"`      // 랭크상태 			set.ranked
	Nsfw       string `query:"nsfw" json:"nsfw"`     // R18				set.nsfw
	Video      string `query:"v" json:"video"`       // 비디오				set.video
	Storyboard string `query:"sb" json:"storyboard"` // 스토리보드			set.storyboard
	//Creator    string `query:"creator" json:"creator"` // 제작자				set.creator

	// map
	Mode             string `query:"m" json:"m"`      // 게임모드			map.mode_int
	TotalLength      minMax `json:"totalLength"`      // 플레이시간			map.totalLength
	MaxCombo         minMax `json:"maxCombo"`         // 콤보				map.maxCombo
	DifficultyRating minMax `json:"difficultyRating"` // 난이도				map.difficultyRating
	Accuracy         minMax `json:"accuracy"`         // od					map.accuracy
	Ar               minMax `json:"ar"`               // ar					map.ar
	Cs               minMax `json:"cs"`               // cs					map.cs
	Drain            minMax `json:"drain"`            // hp					map.drain
	Bpm              minMax `json:"bpm"`              // bpm				map.bpm

	// query
	Sort       string   `query:"sort" json:"sort"`   // 정렬	  order by
	Page       string   `query:"p" json:"page"`      // 페이지 limit
	PageSize   string   `query:"ps" json:"pageSize"` // 페이지 당 크기
	Text       string   `query:"q" json:"query"`     // 문자열 검색
	ParsedText []string `json:"-"`                   // 문자열 검색 파싱
	Option     string   `query:"option" json:"option"`
	OptionB    byte     //artist 1,creator 2,tags 4 ,title 8
	B64        string   `query:"b64"` // body

	//etc
	MapSetId int `param:"si"` // 맵셋id로 검색
	MapId    int `param:"mi"` // 맵id로 검색
}

var searchBaseQuery = `
SELECT 
	beatmapset_id, artist, artist_unicode, creator, favourite_count,
	hype_current, hype_required, nsfw, play_count, source, status,
	title, title_unicode, user_id, video, availability_download_disabled,
	availability_more_information, bpm, can_be_hyped, discussion_enabled,
	discussion_locked, is_scoreable, last_updated, legacy_thread_url,
	nominations_summary_current, nominations_summary_required, ranked,
	ranked_date, storyboard, submitted_date, tags, has_favourited,
	description, genre_id, genre_name, language_id, language_name, ratings
from `

func (s *SearchQuery) queryBuilder2() (qs string, args []interface{}) {
	s.parseQuery()

	var query bytes.Buffer
	var setAnd []string // 맵셋 	AND 문
	var mapAnd []string // 맵	AND 문
	var textSearchQuery []string
	if s.Text != "" {
		text := splitString(s.Text)
		text = utils.MakeArrayUnique(&text)

		if s.OptionB&0x01 == 0x01 {
			textSearchQuery = append(textSearchQuery,
				`SELECT BEATMAPSET_ID from SEARCH_CACHE_ARTIST  WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`,
			)
		}

		if s.OptionB&0x02 == 0x02 {
			textSearchQuery = append(textSearchQuery,
				`SELECT BEATMAPSET_ID from SEARCH_CACHE_CREATOR WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`,
			)
		}

		if s.OptionB&0x04 == 0x04 {
			textSearchQuery = append(textSearchQuery,
				`SELECT BEATMAPSET_ID from SEARCH_CACHE_TAG     WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`,
			)
		}

		if s.OptionB&0x08 == 0x08 {
			textSearchQuery = append(textSearchQuery,
				`SELECT BEATMAPSET_ID from SEARCH_CACHE_TITLE   WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`,
			)
		}

		if s.OptionB == 0xFF {
			textSearchQuery = append(textSearchQuery,
				`SELECT BEATMAPSET_ID from SEARCH_CACHE_OTHER   WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`,
			)
		}

		args = append(args, sql.Named("text", text))           //TODO 검색어 어레이
		args = append(args, sql.Named("textCount", len(text))) //TODO 검색어 어레이.len
		query.WriteString(`WITH SCSI AS (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @text) `)

	}
	if s.Ranked != "all" {
		setAnd = append(setAnd, "ranked IN @ranked")
		args = append(args, sql.Named("ranked", utils.TernaryOperator(ranked[s.Ranked] != nil, ranked[s.Ranked], ranked["default"])))
	}
	if s.Nsfw != "all" {
		setAnd = append(setAnd, "nsfw = "+s.Nsfw)
	}
	if s.Video != "all" {
		setAnd = append(setAnd, "video = "+s.Video)
	}
	if s.Storyboard != "all" {
		setAnd = append(setAnd, "storyboard = "+s.Storyboard)
	}

	if s.Mode != "all" && s.Mode != "" {

		mapAnd = append(mapAnd, "mode_int IN (@modes)")
		args = append(args, sql.Named("modes", mode[s.Mode]))
	}
	if !s.TotalLength.minMaxIsNil() {
		mapAnd = append(mapAnd, `total_length `+s.TotalLength.minMaxAsQuery())
	}
	if !s.MaxCombo.minMaxIsNil() {
		mapAnd = append(mapAnd, `max_combo `+s.MaxCombo.minMaxAsQuery())
	}
	if !s.DifficultyRating.minMaxIsNil() {
		mapAnd = append(mapAnd, `difficulty_rating `+s.DifficultyRating.minMaxAsQuery())
	}
	if !s.Accuracy.minMaxIsNil() {
		mapAnd = append(mapAnd, `accuracy `+s.Accuracy.minMaxAsQuery())
	}
	if !s.Ar.minMaxIsNil() {
		mapAnd = append(mapAnd, `ar `+s.Ar.minMaxAsQuery())
	}
	if !s.Cs.minMaxIsNil() {
		mapAnd = append(mapAnd, `cs `+s.Cs.minMaxAsQuery())
	}
	if !s.Drain.minMaxIsNil() {
		mapAnd = append(mapAnd, `drain `+s.Drain.minMaxAsQuery())
	}
	if !s.Bpm.minMaxIsNil() {
		mapAnd = append(mapAnd, `bpm `+s.Bpm.minMaxAsQuery())
	}

	if len(mapAnd) > 0 { // beatmapset_id IN ()
		mapAnd = append([]string{"beatmapset_id IN (A.beatmapset_id)"}, mapAnd...)
		setAnd = append(setAnd, "beatmapset_id IN (select beatmapset_id from "+config.Config.Sql.Table.Beatmap+" where "+strings.Join(mapAnd, " AND ")+" )")
	}

	query.WriteString(searchBaseQuery)

	query.WriteString(config.Config.Sql.Table.BeatmapSet + " AS MAPSET ")
	if len(textSearchQuery) > 0 {
		query.WriteString(" INNER JOIN ( SELECT * FROM (")
		query.WriteString(strings.Join(textSearchQuery, " UNION ALL "))
		query.WriteString(")A GROUP BY BEATMAPSET_ID having count(*) >= @textCount) AS TEXTINDEX on MAPSET.beatmapset_id = TEXTINDEX.BEATMAPSET_ID")
	}
	if len(setAnd) > 0 { // SELECT * FROM osu.beatmapset WHERE ranked in (4,2,1) AND nsfw = 1 ...
		query.WriteString(" WHERE " + strings.Join(setAnd, " AND "))
	}

	query.WriteString(" ORDER BY " + utils.TernaryOperator(orderBy[s.Sort] == "", orderBy["default"], orderBy[s.Sort]) + " " + s.Page + ";")
	qs = query.String()

	return

}

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(regexpReplace.ReplaceAllString(input, " ")), " ") {
		if s == "" || s == " " {
			continue
		}
		//ss = append(ss, s, porter2.Stemmer.Stem(s))
		ss = append(ss, s)
	}
	return
}
func Search(c echo.Context) (err error) {

	var sq SearchQuery

	err = c.Bind(&sq)
	if sq.B64 != "" {
		b6, err := base64.StdEncoding.DecodeString(sq.B64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-000-0",
				Error:   err,
				Message: "request parm 'b64' base64 decode fail.",
			}))
		}
		err = json.Unmarshal(b6, &sq)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-000-1",
				Error:   err,
				Message: "request parm 'b64' json parse fail.",
			}))
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:    "SEARCH-001",
			Error:   err,
			Message: "request parse error",
		}))
	}

	q, args := sq.queryBuilder2()

	pterm.Info.Println(q)
	rows, err := db.Gorm.Raw(q, args...).Rows()

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-003",
				Error:   err,
				Message: "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:    "SEARCH-004",
			Error:   err,
			Message: "database Query error",
		}))

	}
	defer rows.Close()
	var sets []osu.BeatmapSetsOUT
	var index = map[int]int{}
	var mapids []int

	for rows.Next() {
		var set osu.BeatmapSetsOUT

		err = rows.Scan(
			// beatmapset_id, artist, artist_unicode, creator, favourite_count, hype_current,
			//hype_required, nsfw, play_count, source, status, title, title_unicode, user_id,
			//video, availability_download_disabled, availability_more_information, bpm, can_be_hyped,
			//discussion_enabled, discussion_locked, is_scoreable, last_updated, legacy_thread_url,
			//nominations_summary_current, nominations_summary_required, ranked, ranked_date, storyboard,
			//submitted_date, tags, has_favourited, description, genre_id, genre_name, language_id, language_name, ratings

			&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current, &set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId, &set.Video, &set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped, &set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl, &set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard, &set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name, &set.Language.Id, &set.Language.Name, &set.RatingsString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-005",
				Error:   err,
				Message: "database Query scan error",
			}))
		}

		lu, err := time.Parse("2006-01-02 15:04:05", *set.LastUpdated)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-005-1",
				Error:   err,
				Message: "time Parse error",
			}))
		}
		//pterm.Info.Println(*set.Id, src.FileList[*set.Id].Unix() >= lu.Unix())
		//pterm.Info.Println((*set.Id)*-1, src.FileList[(*set.Id)*-1].Unix() >= lu.Unix())
		set.Cache.Video = src.FileList[*set.Id].Unix() >= lu.Unix()
		set.Cache.NoVideo = src.FileList[(*set.Id)*-1].Unix() >= lu.Unix()

		index[*set.Id] = len(sets)
		mapids = append(mapids, *set.Id)
		sets = append(sets, set)

	}

	if len(sets) < 1 {
		return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:    "SEARCH-006",
			Error:   errors.New(http.StatusText(http.StatusNotFound)),
			Message: "not in database",
		}))

	}

	rows, err = db.Maria.Query(fmt.Sprintf("select beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating, version, accuracy, ar, cs, drain, bpm, "+
		"`convert`, "+
		"count_circles, count_sliders, count_spinners, deleted_at, hit_length, is_scoreable, last_updated, passcount, playcount, checksum, "+
		"user_id from osu.beatmap where beatmapset_id in( %s ) order by difficulty_rating;", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]")))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "SEARCH-007",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err,
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:    "SEARCH-008",
			Error:   err,
			Message: "database Query error",
		}))
	}
	defer rows.Close()

	for rows.Next() {
		var Map osu.BeatmapOUT
		err = rows.Scan(
			//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
			//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
			//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
			&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating, &Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt, &Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:    "SEARCH-009",
				Error:   err,
				Message: "database Query scan error",
			}))
		}
		sets[index[*Map.BeatmapsetId]].Beatmaps = append(sets[index[*Map.BeatmapsetId]].Beatmaps, Map)
	}

	return c.JSON(http.StatusOK, sets)

}
