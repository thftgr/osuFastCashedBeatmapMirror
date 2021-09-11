package Route

import (
	"bytes"
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


	noVideo, err := strconv.ParseBool(c.QueryParam("noVideo")) //1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
	if err != nil {
		noVideo = false
	}

	//================ DEV
	noVideo2, err := strconv.ParseBool(c.QueryParam("nv"))
	if err != nil {
		noVideo2 = false
	}
	noVideo = noVideo || noVideo2
	//================ DEV

	mid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}

	go src.ManualUpdateBeatmapSet(mid)

	rows, err := src.Maria.Query(`SELECT beatmapset_id,artist,title,last_updated,video FROM osu.beatmapset WHERE beatmapset_id = ?`, mid)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		return c.String(404, "Please try again in a few seconds. OR map is not alive. check beatmapset id.")
	}
	var a struct {
		Id          string
		Artist      string
		Title       string
		LastUpdated string
		Video       bool
	}
	if err = rows.Scan(&a.Id, &a.Artist, &a.Title, &a.LastUpdated, &a.Video); err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}

	lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
	if err != nil {
		c.NoContent(http.StatusInternalServerError)
		return
	}
	c.Response().Header().Set("Content-Type", "application/x-osu-beatmap-archive")
	//c.Response().Header().Set("Content-Type", "application/download")

	var serverFileName string
	var url string
	if a.Video && noVideo {
		serverFileName = fmt.Sprintf("%s/-%d.osz", src.Setting.TargetDir, mid)
		if src.FileList[mid*(-1)].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s [no video].osz", a.Id, a.Artist, a.Title))
		}
		url = fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download?noVideo=1", mid)
	} else {
		serverFileName = fmt.Sprintf("%s/%d.osz", src.Setting.TargetDir, mid)
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s.osz", a.Id, a.Artist, a.Title))
		}
		url = fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download", mid)

	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	//noVideo=1

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

	if res.StatusCode != 200 {
		c.NoContent(http.StatusNotFound)
		return
	}
	pterm.Info.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", res.Header.Get("Content-Disposition"))

	var buf bytes.Buffer // 서버에 저장할 파일 버퍼
	//TODO https 응답 먼저 주고 file 저장은 버퍼로 진행

	for i := 0; i < cLen; { // 읽을 데이터 사이즈 체크
		var b = make([]byte, 64000) // 바이트 배열
		n, err := res.Body.Read(b)  // 반쵸 스트림에서 64k 읽어서 바이트 배열 b 에 넣음

		i += n           // 현재까지 읽은 바이트
		buf.Write(b[:n]) // 서버에 저장할 파일 버퍼에 쓴다

		if _, err := c.Response().Write(b[:n]); err != nil { // 클라이언트 리스폰스 스트림에 쓴다(클라이언트 버퍼라 보면 댐)
			c.NoContent(http.StatusInternalServerError) //에러처리
			return err                                  //에러처리
		}
		if err == io.EOF { //에러처리
			break
		} else if err != nil { //에러처리
			fmt.Println(err.Error())
			break
		}
	}
	c.Response().Flush()
	if a.Video && noVideo {
		return saveLocal(&buf, serverFileName, mid*(-1))
	}
	return saveLocal(&buf, serverFileName, mid) // 서버에 파일버퍼를 쓴다

}
func saveLocal(data *bytes.Buffer, path string, id int) (err error) {

	file, err := os.Create(path + ".down")
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
	err = os.Rename(path+".down", path)
	if err != nil {
		return
	}

	src.FileList[id] = time.Now()
	pterm.Info.Println("beatmapSet Downloading Finished", path)
	return
}
