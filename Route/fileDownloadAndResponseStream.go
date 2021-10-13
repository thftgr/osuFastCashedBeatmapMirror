package Route

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func DownloadBeatmapSet(c echo.Context) (err error) {

	//1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
	noVideo, _ := strconv.ParseBool(c.QueryParam("noVideo"))
	noVideo2, _ := strconv.ParseBool(c.QueryParam("nv"))
	noVideo = noVideo || noVideo2

	mid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}

	go src.ManualUpdateBeatmapSet(mid)

	row := src.Maria.QueryRow(`SELECT beatmapset_id,artist,title,last_updated,video FROM osu.beatmapset WHERE beatmapset_id = ?`, mid)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			c.String(http.StatusNotFound, "Please try again in a few seconds. OR map is not alive. check beatmapset id.")
		}
		c.NoContent(http.StatusInternalServerError)
		return
	}

	var a struct {
		Id          string
		Artist      string
		Title       string
		LastUpdated string
		Video       bool
	}

	if err = row.Scan(&a.Id, &a.Artist, &a.Title, &a.LastUpdated, &a.Video); err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}

	lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	c.Response().Header().Set("Content-Type", "application/x-osu-beatmap-archive")

	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download", mid)
	if a.Video && noVideo {
		mid *= -1
		a.Title += " [no video]"
		url += "?noVideo=1"
	}

	serverFileName := fmt.Sprintf("%s/%d.osz", src.Setting.TargetDir, mid)

	if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
		return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s.osz", a.Id, a.Artist, a.Title))
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", src.Setting.Osu.Token.TokenType+" "+src.Setting.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.NoContent(http.StatusNotFound)
		return
	}
	pterm.Info.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", res.Header.Get("Content-Disposition"))

	var buf bytes.Buffer

	//http://localhost/d/1573058
	//http://localhost/d/1469677

	defer c.Response().Flush()

	for i := 0; i < cLen; { // 읽을 데이터 사이즈 체크
		var b = make([]byte, 64000) // 바이트 배열
		n, err := res.Body.Read(b)  // 반쵸 스트림에서 64k 읽어서 바이트 배열 b 에 넣음


		i += n           // 현재까지 읽은 바이트
		if n > 0 {
			buf.Write(b[:n]) // 서버에 저장할 파일 버퍼에 쓴다
			if _, err := c.Response().Write(b[:n]); err != nil {
				c.NoContent(http.StatusInternalServerError)
				return err
			}
		}

		if err == io.EOF {
			break
		} else if err != nil { //에러처리
			fmt.Println(err.Error())
			return err
		}
	}
	if cLen == buf.Len() {
		return saveLocal(&buf, serverFileName, mid)
	}
	errMsg := fmt.Sprintf("filesize not match: bancho response bytes : %d | downloaded bytes : %d",cLen,buf.Len())
	pterm.Error.Printfln(errMsg)
	return errors.New(errMsg)

}
func saveLocal(data *bytes.Buffer, path string, id int) (err error) {
	tmp := path + ".down"
	file, err := os.Create(tmp)
	if err != nil {
		return
	}
	if file == nil {
		return errors.New("")
	}
	_, err = file.Write(data.Bytes())
	if err != nil {
		return
	}
	file.Close()

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	err = os.Rename(tmp, path)
	if err != nil {
		return
	}

	src.FileList[id] = time.Now()
	pterm.Info.Println("beatmapSet Downloading Finished", path)
	return
}
