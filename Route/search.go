package Route

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//
func (s *SearchQuery) parseSort() { //sort

	ss := strings.ToLower(s.Sort)
	switch ss {
	case "ranked_asc", "ranked_date ASC":
		s.Sort = "ranked_date ASC"

	case "ranked_desc", "ranked_date DESC":
		s.Sort = "ranked_date DESC"

	case "favourites_asc", "favourite_count ASC":
		s.Sort = "favourite_count ASC"

	case "favourites_desc", "favourite_count DESC":
		s.Sort = "favourite_count DESC"

	case "plays_asc", "play_count ASC":
		s.Sort = "play_count ASC"

	case "plays_desc", "play_count DESC":
		s.Sort = "play_count DESC"

	case "updated_asc", "last_updated ASC":
		s.Sort = "last_updated ASC"

	case "updated_desc", "last_updated DESC":
		s.Sort = "last_updated DESC"

	case "title_desc", "title DESC":
		s.Sort = "title DESC"

	case "title_asc", "title ASC":
		s.Sort = "title ASC"

	case "artist_desc", "artist DESC":
		s.Sort = "artist DESC"

	case "artist_asc", "artist ASC":
		s.Sort = "artist ASC"
	default:
		s.Sort = "ranked_date DESC"
	}
}

func (s *SearchQuery) parsePage() {
	atoi, err := strconv.Atoi(s.Page)
	if err != nil || atoi <= 0 {
		s.Page = "LIMIT 50"
	} else {
		s.Page = fmt.Sprintf("LIMIT %d,50", atoi*50)
	}
	return
}

func (s *SearchQuery) parseMode() {
	switch s.Mode {
	case "0":
		s.Mode = "0"
	case "1":
		s.Mode = "1"
	case "2":
		s.Mode = "2"
	case "3":
		s.Mode = "3"
	default:
		s.Mode = "all"
	}
}

func (s *SearchQuery) parseRanked() {
	ss := strings.ToLower(s.Ranked)
	switch ss {
	case "ranked", "1,2":
		s.Ranked = "1,2"

	case "qualified":
		s.Ranked = "3"

	case "loved":
		s.Ranked = "4"

	case "pending":
		s.Ranked = "0"

	case "wip":
		s.Ranked = "-1"

	case "graveyard":
		s.Ranked = "-2"

	case "unranked", "0,-1,-2":
		s.Ranked = "0,-1,-2"

	case "any":
		s.Ranked = "all"

	case "-2", "-1", "0", "1", "2", "3", "4":
		s.Ranked = ss

	default:
		s.Ranked = "4,2,1"

	}
}
func (s *SearchQuery) parseNsfw() {
	ss := strings.ToLower(s.Ranked)
	switch ss {
	case "1":
		s.Nsfw = "all"
	default:
		s.Nsfw = "0"
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
		return fmt.Sprintf("> %.1f", v.Min)
	} else if v.Min == 0 && v.Max != 0 {
		return fmt.Sprintf("< %.1f", v.Max)
	}
	return fmt.Sprintf("BETWEEN %.1f AND %.1f", v.Min, v.Max)

}
func (s *SearchQuery) parseQuery() {
	s.parseRanked()
	s.parseMode()
	s.parseSort()
	s.parseNsfw()
	s.parsePage()
	s.parseExtra()
}

type SearchQuery struct {
	// global
	Extra string `query:"e" json:"extra"` // 스토리보드 비디오.

	// set
	Ranked     string `query:"s" json:"ranked"`        // 랭크상태 			set.ranked
	Nsfw       string `query:"nsfw" json:"nsfw"`       // R18				set.nsfw
	Video      string `query:"v" json:"video"`         // 비디오				set.video
	Storyboard string `query:"sb" json:"storyboard"`   // 스토리보드			set.storyboard
	Creator    string `query:"creator" json:"creator"` // 제작자				set.creator

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
	Sort string `query:"sort" json:"sort"` // 정렬				order by
	Page string `query:"p" json:"page"`    // 페이지				limit
	Text string `query:"q" json:"query"`   // 문자열 검색

	//etc
	MapSetId int `param:"si"` // 맵셋id로 검색
	MapId    int `param:"mi"` // 맵id로 검색
}

func queryBuilder(s *SearchQuery) (qs string, i []interface{}) {
	s.parseQuery()

	var query bytes.Buffer
	var buf1 []string
	var buf2 []string
	query.WriteString("SELECT * FROM osu.beatmapset ")
	//ranked

	if s.Ranked != "all" {
		buf1 = append(buf1, "RANKED IN("+s.Ranked+")")
		buf2 = append(buf2, "RANKED IN("+s.Ranked+")")
	}

	if s.Mode != "all" {
		buf2 = append(buf2, "mode_int IN("+s.Mode+")")
	}

	if len(buf2) > 0 {
		buf1 = append(buf1, "beatmapset_id IN (SELECT DISTINCT beatmapset_id FROM osu.beatmap WHERE "+strings.Join(buf2, " AND ")+" )")
	}

	if s.Text != "" {
		buf1 = append(buf1, "beatmapset_id IN (SELECT beatmapset_id FROM osu.search_index WHERE MATCH(text) AGAINST(?))")
		i = append(i, s.Text)

	}

	if len(buf1) > 0 {
		query.WriteString("WHERE ")
		query.WriteString(strings.Join(buf1, " AND "))
	}
	query.WriteString("ORDER BY ")
	query.WriteString(s.Sort)
	query.WriteString(" ")
	query.WriteString(s.Page)
	query.WriteString(";")
	qs = query.String()

	return

}

func Search(c echo.Context) (err error) {

	var sq SearchQuery
	err = c.Bind(&sq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Logger.Error(&bodyStruct.ErrorStruct{
			Code:      "SEARCH-001",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err.Error(),
			Message:   "request parse error",
		}))

	}

	q, qs := queryBuilder(&sq)
	go func() {
		b, _ := json.Marshal(sq)
		log.Println("REBUILDED REQUEST:", string(b))
		log.Println("GENERATED QUERY:", q, "ARGS:", qs)
		t := time.Now().Format("2006/01/02 15:01:05") //2021/09/10 22:30:38
		pterm.Info.Println(t, "REBUILDED REQUEST:", pterm.LightYellow(string(b)))
		pterm.Info.Println(t, "GENERATED QUERY:", pterm.LightYellow(q), "ARGS:", pterm.LightYellow(qs))
	}()

	rows, err := db.Maria.Query(q, qs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, Logger.Error(&bodyStruct.ErrorStruct{
				Code:      "SEARCH-002",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err.Error(),
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, Logger.Error(&bodyStruct.ErrorStruct{
			Code:      "SEARCH-003",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err.Error(),
			Message:   "database Query error",
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
			return c.JSON(http.StatusInternalServerError, Logger.Error(&bodyStruct.ErrorStruct{
				Code:      "SEARCH-005",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err.Error(),
				Message:   "database Query scan error",
			}))
		}
		index[*set.Id] = len(sets)
		mapids = append(mapids, *set.Id)
		sets = append(sets, set)
	}

	if len(sets) < 1 {
		return c.JSON(http.StatusNotFound, Logger.Error(&bodyStruct.ErrorStruct{
			Code:      "SEARCH-006",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     http.StatusText(http.StatusNotFound),
			Message:   "not in database",
		}))

	}

	rows, err = db.Maria.Query(fmt.Sprintf(`select * from osu.beatmap where beatmapset_id in( %s ) order by difficulty_rating asc;`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]")))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, Logger.Error(&bodyStruct.ErrorStruct{
				Code:      "SEARCH-007",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err.Error(),
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, Logger.Error(&bodyStruct.ErrorStruct{
			Code:      "SEARCH-008",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err.Error(),
			Message:   "database Query error",
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
			return c.JSON(http.StatusInternalServerError, Logger.Error(&bodyStruct.ErrorStruct{
				Code:      "SEARCH-009",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err.Error(),
				Message:   "database Query scan error",
			}))
		}
		sets[index[*Map.BeatmapsetId]].Beatmaps = append(sets[index[*Map.BeatmapsetId]].Beatmaps, Map)

	}

	return c.JSON(http.StatusOK, sets)

}
