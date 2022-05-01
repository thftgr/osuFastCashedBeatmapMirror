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
)

func SearchByBeatmapId(c echo.Context) (err error) {
	var sq SearchQuery
	err = c.Bind(&sq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapId-001",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "request parse error",
		}))
	}
	fmt.Println(sq.MapId)
	row := db.Maria.QueryRow("SELECT BEATMAP_ID, BEATMAPSET_ID, MODE, MODE_INT, STATUS, RANKED, TOTAL_LENGTH, MAX_COMBO, DIFFICULTY_RATING, VERSION, ACCURACY, AR, CS, DRAIN, BPM, `CONVERT`, COUNT_CIRCLES, COUNT_SLIDERS, COUNT_SPINNERS, DELETED_AT, HIT_LENGTH, IS_SCOREABLE, LAST_UPDATED, PASSCOUNT, PLAYCOUNT, CHECKSUM, USER_ID FROM BEATMAP WHERE BEATMAP_ID = ?;", sq.MapId)
	var Map osu.BeatmapOUT
	err = row.Scan(
		//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
		//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
		//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
		&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating,
		&Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert, &Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt,
		&Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "SearchByBeatmapId-002",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     err,
				Message:   "not in database",
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "SearchByBeatmapId-003",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "database Query error",
		}))
	}

	return c.JSON(http.StatusOK, Map)
}
