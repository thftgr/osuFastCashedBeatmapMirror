package src

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type FileIndex map[int]time.Time

var FileList = make(FileIndex)

func StartIndex() {
	for {
		FileListUpdate()
		time.Sleep(time.Second * 60 * 5)
	}
}

func FileListUpdate() {

	files, err := ioutil.ReadDir(Setting.TargetDir)
	if err != nil {
		panic(err)
	}

	tmp := make(FileIndex)
	for _, file := range files {
		if sid, err := strconv.Atoi(strings.Replace(file.Name(), ".osz", "", -1)); err == nil {
			tmp[sid] = file.ModTime()
		}
	}
	FileList = tmp
	fmt.Printf(
		"%s File List Indexing : %s files\n",
		time.Now().Format("2006-01-02 15:04:05"),
		[10]string{strconv.Itoa(len(FileList))},
	)

}
