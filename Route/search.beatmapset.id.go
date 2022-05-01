package Route

import (
	"database/sql"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func SearchByBeatmapSetId(c echo.Context) (err error) {
	var sq SearchQuery
	err = c.Bind(&sq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-001",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "request parse error",
		}))
	}
	row := db.Maria.QueryRow(`SELECT 
       BEATMAPSET_ID, ARTIST, ARTIST_UNICODE, CREATOR, FAVOURITE_COUNT, 
       HYPE_CURRENT, HYPE_REQUIRED, NSFW, PLAY_COUNT, SOURCE, STATUS, 
       TITLE, TITLE_UNICODE, USER_ID, VIDEO, AVAILABILITY_DOWNLOAD_DISABLED, 
       AVAILABILITY_MORE_INFORMATION, BPM, CAN_BE_HYPED, DISCUSSION_ENABLED, 
       DISCUSSION_LOCKED, IS_SCOREABLE, LAST_UPDATED, LEGACY_THREAD_URL, 
       NOMINATIONS_SUMMARY_CURRENT, NOMINATIONS_SUMMARY_REQUIRED, RANKED, 
       RANKED_DATE, STORYBOARD, SUBMITTED_DATE, TAGS, HAS_FAVOURITED, DESCRIPTION, 
       GENRE_ID, GENRE_NAME, LANGUAGE_ID, LANGUAGE_NAME, RATINGS FROM BEATMAPSET WHERE BEATMAPSET_ID = ?;`, sq.MapSetId)

	var set osu.BeatmapSetsOUT
	var mapids []int
	err = row.Scan(
		// beatmapset_id, artist, artist_unicode, creator, favourite_count, hype_current,
		//hype_required, nsfw, play_count, source, status, title, title_unicode, user_id,
		//video, availability_download_disabled, availability_more_information, bpm, can_be_hyped,
		//discussion_enabled, discussion_locked, is_scoreable, last_updated, legacy_thread_url,
		//nominations_summary_current, nominations_summary_required, ranked, ranked_date, storyboard,
		//submitted_date, tags, has_favourited, description, genre_id, genre_name, language_id, language_name, ratings

		&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current, &set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId, &set.Video, &set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped, &set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl, &set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard, &set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name, &set.Language.Id, &set.Language.Name, &set.RatingsString)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-002",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err,
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-003",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "database Query error",
		}))
	}
	mapids = append(mapids, *set.Id)

	if *set.Id == 0 {
		return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-004",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "not in database",
		}))
	}

	rows, err := db.Maria.Query(fmt.Sprintf(
		"SELECT BEATMAP_ID, BEATMAPSET_ID, MODE, MODE_INT, STATUS, RANKED, TOTAL_LENGTH, MAX_COMBO, DIFFICULTY_RATING, VERSION, ACCURACY, AR, CS, DRAIN, BPM, `CONVERT`, COUNT_CIRCLES, "+
			"COUNT_SLIDERS, COUNT_SPINNERS, DELETED_AT, HIT_LENGTH, IS_SCOREABLE, LAST_UPDATED, PASSCOUNT, PLAYCOUNT, CHECKSUM, "+
			"USER_ID FROM BEATMAP WHERE BEATMAPSET_ID IN( %s ) ORDER BY DIFFICULTY_RATING;", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(mapids)), ", "), "[]")))

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-005",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err,
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapSetId-006",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
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
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapSetId-007",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err,
				Message:   "database Query scan error",
			}))
		}
		set.Beatmaps = append(set.Beatmaps, Map)

	}

	return c.JSON(http.StatusOK, set)
}
