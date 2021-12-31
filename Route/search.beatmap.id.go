package Route

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/osu"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"net/http"
)

func SearchByBeatmapId(c echo.Context) (err error) {
	var sq SearchQuery
	err = c.Bind(&sq)
	if err != nil {
		pterm.Error.Println(err)
		c.NoContent(http.StatusInternalServerError)
		return
	}
	fmt.Println(sq.MapId)
	row := src.Maria.QueryRow(`select * from osu.beatmap where beatmap_id = ?;`, sq.MapId)
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
		pterm.Error.Println(err)
		c.NoContent(http.StatusInternalServerError)
		return
	}

	return c.JSON(http.StatusOK, Map)
}
