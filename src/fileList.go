package src

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FileIndex map[int]time.Time

var FileList = make(FileIndex)
var fileSize uint64

const goos = runtime.GOOS

func StartIndex() {
	FileListUpdate()
	go func() {
		time.Sleep(time.Second * 60 * 5)
		for {
			FileListUpdate()
			time.Sleep(time.Second * 60 * 5)
		}
	}()

}
func FileListUpdate() {
	var err error

	checkDir()
	files, err := ioutil.ReadDir(config.Setting.TargetDir)
	if err != nil {
		return
	}

	tmp := make(FileIndex)
	fileSize = 0
	for _, file := range files {
		if sid, err := strconv.Atoi(strings.Replace(file.Name(), ".osz", "", -1)); err == nil {
			tmp[sid] = file.ModTime()
			fileSize += uint64(file.Size())
		}
	}
	FileList = tmp
	pterm.Info.Printfln(
		"%s File List Indexing : %s files [%s]\n",
		time.Now().Format("2006-01-02 15:04:05"),
		pterm.LightYellow(strconv.Itoa(len(FileList))),
		pterm.LightYellow(totalFileSize()),
	)

}

func totalFileSize() (s string) {
	if goos == "windows" {
		if fileSize > 1099511627776 { //TB
			return fmt.Sprintf("%d%s", fileSize/1099511627776, "TB")
		} else if fileSize > 1073741824 { //GB
			return fmt.Sprintf("%d%s", fileSize/1073741824, "GB")
		} else if fileSize > 1048576 { //MB
			return fmt.Sprintf("%d%s", fileSize/1048576, "MB")
		} else if fileSize > 1024 { //KB
			return fmt.Sprintf("%d%s", fileSize/1024, "KB")
		}
	} else {
		if fileSize > 1000000000000 { //TB
			return fmt.Sprintf("%d%s", fileSize/1000000000000, "TB")
		} else if fileSize > 1000000000 { //GB
			return fmt.Sprintf("%d%s", fileSize/1000000000, "GB")
		} else if fileSize > 1000000 { //MB
			return fmt.Sprintf("%d%s", fileSize/1000000, "MB")
		} else if fileSize > 1000 { //KB
			return fmt.Sprintf("%d%s", fileSize/1000, "KB")
		}
	}

	return fmt.Sprintf("%d%s", fileSize, "B")
}
func checkDir() {
	if _, e := os.Stat(config.Setting.TargetDir); os.IsNotExist(e) {
		err := os.MkdirAll(config.Setting.TargetDir, 666)
		if err != nil {
			pterm.Error.Println(err)
			panic(err)
		}
	}
}
