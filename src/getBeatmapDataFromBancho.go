package src

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/osu"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var api = struct {
	count int
	mutex sync.Mutex
}{}
var pause bool

func apicountAdd() {
	api.mutex.Lock()
	api.count++
	api.mutex.Unlock()
}

func apiCountReset() {
	api.mutex.Lock()
	api.count = 0
	api.mutex.Unlock()
}
func RunGetBeatmapDataASBancho() {
	go func() {
		for {
			time.Sleep(time.Minute)
			if Maria.Ping() != nil {
				continue
			}
			apiCountReset()
			go Setting.Save()
		}
	}()
	go func() { //ALL desc limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Second * 30)
			getUpdatedMapDesc()
		}
	}()
	go func() { //Update Ranked DESC limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getUpdatedMapRanked()
		}
	}()
	go func() { //Update Loved DESC limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getUpdatedMapLoved()
		}
	}()
	go func() { //Update Graveyard asc limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getGraveyardMap()
		}
	}()
	go func() { //ALL asc
		for {
			awaitApiCount()
			getUpdatedMapAsc()
		}
	}()
	pterm.Info.Println("Bancho cron started.")
}
func awaitApiCount() {
	for {
		if api.count < 60 && !pause {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}
func ManualUpdateBeatmapSet(id int) {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d", id)

	var data osu.BeatmapSetsIN
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	updateMapset(&data)
}

func getUpdatedMapRanked() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

	var data osu.BeatmapsetsSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}

	return
}

func getUpdatedMapLoved() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=loved"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	return
}

func getGraveyardMap() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := ""
	lu := &Setting.Osu.BeatmapUpdate.GraveyardAsc.LastUpdate
	id := &Setting.Osu.BeatmapUpdate.GraveyardAsc.Id
	if *lu+*id != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard&cursor%5Blast_update%5D=" + *lu + "&cursor%5B_id%5D=" + *id
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}
	if data.Cursor == nil {
		*lu = ""
		*id = ""
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	*lu = *data.Cursor.LastUpdate
	*id = *data.Cursor.Id
	return
}
func getUpdatedMapDesc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

	var data osu.BeatmapsetsSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	Setting.Osu.BeatmapUpdate.UpdatedDesc.LastUpdate = *data.Cursor.LastUpdate
	Setting.Osu.BeatmapUpdate.UpdatedDesc.Id = *data.Cursor.Id

	return
}

func getUpdatedMapAsc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := ""
	lu := &Setting.Osu.BeatmapUpdate.UpdatedAsc.LastUpdate
	id := &Setting.Osu.BeatmapUpdate.UpdatedAsc.Id
	if *lu+*id != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any&cursor%5Blast_update%5D=" + *lu + "&cursor%5B_id%5D=" + *id
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}
	if data.Cursor == nil {
		*lu = ""
		*id = ""
		return
	}
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	*lu = *data.Cursor.LastUpdate
	*id = *data.Cursor.Id
	return
}

func stdGETBancho(url string, str interface{}) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	req.Header.Add("Authorization", Setting.Osu.Token.TokenType+" "+Setting.Osu.Token.AccessToken)

	res, err := client.Do(req)
	apicountAdd()
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			pause = true
		}
		return errors.New(res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &str)

}

func updateMapset(s *osu.BeatmapSetsIN) {
	//	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	//	hype_current,hype_required,nsfw,play_count,source,
	//	status,title,title_unicode,user_id,video,
	//	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	//	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	//	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	//	tags,has_favourited,description,genre_id,genre_name,
	//	language_id,language_name,ratings

	r := *s.Ratings
	_, err := Maria.Exec(UpsertBeatmapSet, s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount,
		s.Hype.Current, s.Hype.Required, s.Nsfw, s.PlayCount, s.Source,
		s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video,
		s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled,
		s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current,
		s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate,
		s.Tags, s.HasFavourited, s.Description.Description, s.Genre.Id, s.Genre.Name,
		s.Language.Id, s.Language.Name, fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10]),
	)
	if err != nil {
		log.Error(err)
		pterm.Error.Println(err)
	}


	if *s.Beatmaps == nil {
		return
	}
	ch := make(chan struct{}, len(*s.Beatmaps))
	for _, m := range *s.Beatmaps {
		go upsertMap(m, ch)
	}
	for i := 0; i < len(ch); i++ {
		<-ch
	}
}

func upsertMap(m osu.BeatmapIN, ch chan struct{}) {
	_, err := Maria.Exec(UpsertBeatmap, m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version,
		m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt,
		m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
	)
	if err != nil {
		log.Error(err)
		pterm.Error.Println(err)
	}

	ch <- struct{}{}
}

const (
	UpsertBeatmap = `
	INSERT INTO osu.beatmap
		(
			beatmap_id,beatmapset_id,mode,mode_int,status,	ranked,total_length,max_combo,difficulty_rating,version,
			accuracy,ar,cs,drain,bpm,` + "`convert`" + `,count_circles,count_sliders,count_spinners,deleted_at,
			hit_length,is_scoreable,last_updated,passcount,playcount,	checksum,user_id
		)VALUES(
			?,?,?,?,?,	?,?,?,?,?,
			?,?,?,?,?,	?,?,?,?,?,
			?,?,?,?,?,	?,?
		)ON DUPLICATE KEY UPDATE 
			beatmapset_id = VALUES(beatmapset_id), mode = VALUES(mode), mode_int = VALUES(mode_int), status = VALUES(status), 
			ranked = VALUES(ranked), total_length = VALUES(total_length), max_combo = VALUES(max_combo), difficulty_rating = VALUES(difficulty_rating), 
			version = VALUES(version), 	accuracy = VALUES(accuracy), ar = VALUES(ar), cs = VALUES(cs), drain = VALUES(drain), bpm = VALUES(bpm),` +
		"`convert` = VALUES(`convert`" + `), count_circles = VALUES(count_circles), count_sliders = VALUES(count_sliders), 
			count_spinners = VALUES(count_spinners), deleted_at = VALUES(deleted_at), 	hit_length = VALUES(hit_length), 
			is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), passcount = VALUES(passcount), playcount = VALUES(playcount), 
			checksum = VALUES(checksum), user_id = VALUES(user_id);`

	setUpsert = `
		INSERT INTO osu.beatmapset (
			beatmapset_id,artist,artist_unicode,creator,favourite_count,
			nsfw,play_count,source,
			status,title,title_unicode,user_id,video,
			availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
			discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
			nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
			tags,has_favourited )
		VALUES %s ON DUPLICATE KEY UPDATE 
			artist = VALUES(artist), artist_unicode = VALUES(artist_unicode), creator = VALUES(creator), favourite_count = VALUES(favourite_count), 
			nsfw = VALUES(nsfw), play_count = VALUES(play_count), source = VALUES(source), 
			status = VALUES(status), title = VALUES(title), title_unicode = VALUES(title_unicode), user_id = VALUES(user_id), video = VALUES(video), 
			availability_download_disabled = VALUES(availability_download_disabled), availability_more_information = VALUES(availability_more_information), 
			bpm = VALUES(bpm), can_be_hyped = VALUES(can_be_hyped), discussion_enabled = VALUES(discussion_enabled), 
			discussion_locked = VALUES(discussion_locked), is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), 
			legacy_thread_url = VALUES(legacy_thread_url), nominations_summary_current = VALUES(nominations_summary_current), 
			nominations_summary_required = VALUES(nominations_summary_required), ranked = VALUES(ranked), ranked_date = VALUES(ranked_date), 
			storyboard = VALUES(storyboard), submitted_date = VALUES(submitted_date), 
			tags = VALUES(tags), has_favourited = VALUES(has_favourited);`
	setValues = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //30

	mapUpsert = `
		INSERT INTO osu.beatmap (	
			beatmap_id,beatmapset_id,mode,mode_int,status,	ranked,total_length,max_combo,difficulty_rating,version,
			accuracy,ar,cs,drain,bpm,` + "`convert`" + `,count_circles,count_sliders,count_spinners,deleted_at,
			hit_length,is_scoreable,last_updated,passcount,playcount,	checksum,user_id
		)VALUES %s ON DUPLICATE KEY UPDATE 
			beatmapset_id = VALUES(beatmapset_id), mode = VALUES(mode), mode_int = VALUES(mode_int), status = VALUES(status), 
			ranked = VALUES(ranked), total_length = VALUES(total_length), max_combo = VALUES(max_combo), 
			difficulty_rating = VALUES(difficulty_rating), version = VALUES(version), 
			accuracy = VALUES(accuracy), ar = VALUES(ar), cs = VALUES(cs), drain = VALUES(drain), bpm = VALUES(bpm), 
			` + "`convert` = VALUES(`convert`" + `), count_circles = VALUES(count_circles), count_sliders = VALUES(count_sliders),
			count_spinners = VALUES(count_spinners), deleted_at = VALUES(deleted_at), 
			hit_length = VALUES(hit_length), is_scoreable = VALUES(is_scoreable), last_updated = VALUES(last_updated), 
			passcount = VALUES(passcount), playcount = VALUES(playcount), 
			checksum = VALUES(checksum), user_id = VALUES(user_id);`
	mapValues         = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //27
	selectDeletedMaps = `select beatmap_id from osu.beatmap where beatmapset_id in (%s) AND beatmap_id not in (%s)`
	deleteMap         = `delete from osu.beatmap where beatmap_id in (%s);`
	UpsertBeatmapSet  = `
INSERT INTO osu.beatmapset(
	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	hype_current,hype_required,nsfw,play_count,source,
	status,title,title_unicode,user_id,video,
	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	tags,has_favourited,description,genre_id,genre_name,
	language_id,language_name,ratings
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?
)ON DUPLICATE KEY UPDATE 
	artist= VALUES(artist), artist_unicode= VALUES(artist_unicode), creator= VALUES(creator), favourite_count= VALUES(favourite_count), 
	hype_current= VALUES(hype_current), hype_required= VALUES(hype_required), nsfw= VALUES(nsfw), play_count= VALUES(play_count), source= VALUES(source), 
	status= VALUES(status), title= VALUES(title), title_unicode= VALUES(title_unicode), user_id= VALUES(user_id), video= VALUES(video), 
	availability_download_disabled= VALUES(availability_download_disabled), availability_more_information= VALUES(availability_more_information), 
	bpm= VALUES(bpm), can_be_hyped= VALUES(can_be_hyped), discussion_enabled= VALUES(discussion_enabled), 
	discussion_locked= VALUES(discussion_locked), is_scoreable= VALUES(is_scoreable), last_updated= VALUES(last_updated), legacy_thread_url= VALUES(legacy_thread_url), 
	nominations_summary_current= VALUES(nominations_summary_current), 	nominations_summary_required= VALUES(nominations_summary_required), 
	ranked= VALUES(ranked), ranked_date= VALUES(ranked_date), storyboard= VALUES(storyboard), submitted_date= VALUES(submitted_date), 
	tags= VALUES(tags), has_favourited= VALUES(has_favourited), description= VALUES(description), genre_id= VALUES(genre_id), genre_name= VALUES(genre_name), 
	language_id= VALUES(language_id), language_name= VALUES(language_name), ratings= VALUES(ratings)
;
`
)

func buildSqlValues(s string, count int) (r string) {
	var sbuf []string
	for i := 0; i < count; i++ {
		sbuf = append(sbuf, s)
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sbuf)), ","), "[]")
}

func updateSearchBeatmaps(data *[]osu.BeatmapSetsIN) (err error) {
	if data == nil {
		return
	}
	if len(*data) < 1 {
		return
	}

	var (
		setInsertBuf []interface{}
		mapInsertBuf []interface{}
		beatmapSets  []int
		beatmaps     []int
		deletedMaps  []int
	)

	for _, s := range *data {
		beatmapSets = append(beatmapSets, s.Id)
		setInsertBuf = append(setInsertBuf,
			s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount, s.Nsfw, s.PlayCount, s.Source,
			s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video,
			s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled,
			s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current,
			s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate,
			s.Tags, s.HasFavourited,
		)
		for _, m := range *s.Beatmaps {
			beatmaps = append(beatmaps, m.Id)
			mapInsertBuf = append(mapInsertBuf,
				m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version,
				m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt,
				m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
			)
		}
	}
	//맵셋
	if _, err = Maria.Exec(fmt.Sprintf(setUpsert, buildSqlValues(setValues, len(beatmapSets))), setInsertBuf...); err != nil {
		pterm.Error.Println(err)
		return err
	}

	//맵
	//fmt.Println(fmt.Sprintf(mapUpsert, buildSqlValues(mapValues, len(beatmaps))))
	if _, err = Maria.Exec(fmt.Sprintf(mapUpsert, buildSqlValues(mapValues, len(beatmaps))), mapInsertBuf...); err != nil {
		pterm.Error.Println(err)
		return err
	}

	//삭제된 맵 제거
	sets := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(beatmapSets)), ","), "[]")
	maps := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(beatmaps)), ","), "[]")
	rows, err := Maria.Query(fmt.Sprintf(selectDeletedMaps, sets, maps))
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var i int
		if err = rows.Scan(&i); err != nil {
			return err
		}
		deletedMaps = append(deletedMaps, i)
	}
	if len(deletedMaps) > 1 {
		pterm.Info.Println(time.Now().Format("02 15:04:05"), "DELETED MAPS:", pterm.LightYellow(deletedMaps))
		dmaps := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(deletedMaps)), ","), "[]")
		if _, err = Maria.Exec(fmt.Sprintf(deleteMap, dmaps)); err != nil {
			return err
		}
	}

	return
}
