package Route

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/thftgr/osuFastCashedBeatmapMirror/db"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

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
	fmt.Println("beatmapSet Downloading Finished", path)
	return
}

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
		c.NoContent(500)
		return
	}

	go src.ManualUpdateBeatmapSet(mid)

	rows, err := db.Maria.Query(db.GetDownloadBeatmapData, mid)
	if err != nil {
		c.NoContent(500)
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
		c.NoContent(500)
		return
	}


	lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
	if err != nil {
		c.NoContent(500)
		return
	}

	var serverFileName string
	var url string
	if a.Video && noVideo {
		serverFileName = fmt.Sprintf("%s/-%d.osz", src.Setting.TargetDir, mid)
		if src.FileList[mid*(-1)].Unix() >= lu.Unix() { // 맵이 최신인경우
			c.Response().Header().Set("Content-Type", "application/download")
			return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s [no video].osz",a.Id,a.Artist ,a.Title))
		}
		url = fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download?noVideo=1",mid)
	} else {
		serverFileName = fmt.Sprintf("%s/%d.osz", src.Setting.TargetDir, mid)
		if src.FileList[mid].Unix() >= lu.Unix() { // 맵이 최신인경우
			c.Response().Header().Set("Content-Type", "application/download")
			return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s.osz",a.Id,a.Artist ,a.Title))
		}
		url = fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download",mid)

	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	//noVideo=1

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		c.NoContent(500)
		return
	}
	req.Header.Add("Authorization", src.Setting.Osu.Token.TokenType+" "+src.Setting.Osu.Token.AccessToken)

	res, err := client.Do(req)

	if err != nil {
		c.NoContent(500)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		c.NoContent(404)
		return
	}

	fmt.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Type", "application/download")
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", res.Header.Get("Content-Disposition"))

	var buf = bytes.Buffer{}
	//TODO https 응답 먼저 주고 file 저장은 버퍼로 진행

	for i := 0; i < cLen; {
		var b = make([]byte, 64000)
		n, err := res.Body.Read(b)

		i += n
		buf.Write(b[:n])

		if _, err := c.Response().Write(b[:n]); err != nil {
			c.NoContent(500)
			return err
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	c.Response().Flush()
	if a.Video && noVideo {
		return saveLocal(&buf, serverFileName, mid*(-1))
	}
	return saveLocal(&buf, serverFileName, mid)

}
