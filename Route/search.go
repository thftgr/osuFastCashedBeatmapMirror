package Route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/osu"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//
func parseSort(s string) (ss string) { //sort

	s = strings.ToLower(s)
	switch s {

	case "ranked_asc":
		ss += "ranked_date ASC"
	case "favourites_asc":
		ss += "favourite_count ASC"
	case "favourites_desc":
		ss += "favourite_count DESC"
	case "plays_asc":
		ss += "play_count ASC"
	case "plays_desc":
		ss += "play_count DESC"
	case "updated_asc":
		ss += "last_updated ASC"
	case "updated_desc":
		ss += "last_updated DESC"
	case "title_desc":
		ss += "title DESC"
	case "title_asc":
		ss += "title ASC"
	case "artist_desc":
		ss += "artist DESC"
	case "artist_asc":
		ss += "artist ASC"
	default:
		ss += "ranked_date DESC"
	}

	return
}

func parsePage(s string) (ss string) {
	atoi, err := strconv.Atoi(s)
	if err != nil || atoi <= 0 {
		return "LIMIT 50"
	}
	return fmt.Sprintf("LIMIT %d,50", atoi*50)
}

func parseMode(s string) (ss string) {
	switch s {
	case "0":
		ss = "0"
	case "1":
		ss = "1"
	case "2":
		ss = "2"
	case "3":
		ss = "3"
	default:
		ss = "all"
	}
	return
}

func parseStatus(s string) (ss string) {
	switch s {
	case "ranked":
		ss = "1,2"
	case "qualified":
		ss = "3"
	case "loved":
		ss = "4"
	case "pending":
		ss = "0"
	case "wip":
		ss = "-1"
	case "graveyard":
		ss = "-2"
	case "any":
		ss = "all"
	default:
		ss = "4,2,1"

	}

	return
}

func queryBuilder(s *SearchQuery) (qs string, i []interface{}) {

	(*s).Status = parseStatus((*s).Status)
	(*s).Mode = parseMode((*s).Mode)
	(*s).Sort = parseSort((*s).Sort)
	(*s).Page = parsePage((*s).Page)

	var query bytes.Buffer
	var buf1 []string
	var buf2 []string
	query.WriteString("SELECT * FROM osu.beatmapset ")
	//ranked

	if (*s).Status != "all" {
		buf1 = append(buf1, "RANKED IN("+(*s).Status+")")
		buf2 = append(buf2, "RANKED IN("+(*s).Status+")")
	}

	if (*s).Mode != "all" {
		buf2 = append(buf2, "mode_int IN("+(*s).Mode+")")
	}

	if len(buf2) > 0 {
		buf1 = append(buf1, "beatmapset_id IN (SELECT DISTINCT beatmapset_id FROM osu.beatmap WHERE "+strings.Join(buf2, " AND ")+" )")
	}

	if (*s).Text != "" {
		buf1 = append(buf1, "beatmapset_id IN (SELECT beatmapset_id FROM osu.search_index WHERE MATCH(text) AGAINST(?))")
		i = append(i, (*s).Text)

	}

	if len(buf1) > 0 {
		query.WriteString("WHERE ")
		query.WriteString(strings.Join(buf1, " AND "))
	}
	query.WriteString("ORDER BY ")
	query.WriteString((*s).Sort)
	query.WriteString(" ")
	query.WriteString((*s).Page)
	query.WriteString(";")
	qs = query.String()

	return

}

type SearchQuery struct {
	Status string `query:"s" json:"s"`
	Mode   string `query:"m" json:"m"`
	Sort   string `query:"sort" json:"sort"`
	Page   string `query:"p" json:"p"`
	Text   string `query:"q" json:"q"`
}

func Search(c echo.Context) (err error) {
	var sq SearchQuery
	err = c.Bind(&sq)
	if err != nil {
		pterm.Error.Println(err)
		c.NoContent(http.StatusInternalServerError)
		return
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

	rows, err := src.Maria.Query(q, qs...)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
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

			&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current,
			&set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId,
			&set.Video, &set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped,
			&set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl,
			&set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard,
			&set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name,
			&set.Language.Id, &set.Language.Name, &set.RatingsString)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		index[*set.Id] = len(sets)
		mapids = append(mapids, *set.Id)
		sets = append(sets, set)
	}

	if len(sets) < 1 {
		c.NoContent(http.StatusNotFound)
		return
	}

	rows, err = src.Maria.Query(fmt.Sprintf(
		`select * from osu.beatmap where beatmapset_id in( %s ) order by difficulty_rating desc;`,
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]"),
	))

	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var Map osu.BeatmapOUT
		err = rows.Scan(
			//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
			//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
			//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
			&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating,
			&Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt,
			&Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId,
		)
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		sets[index[*Map.BeatmapsetId]].Beatmaps = append(sets[index[*Map.BeatmapsetId]].Beatmaps, Map)

	}

	return c.JSON(http.StatusOK, sets)

}
