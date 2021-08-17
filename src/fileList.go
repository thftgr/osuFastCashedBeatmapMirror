package src

import (
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type FileIndex map[int]time.Time

var FileList = make(FileIndex)

func StartIndex(c *pterm.SpinnerPrinter) {
	FileListUpdate(c)
	go func() {
		time.Sleep(time.Second * 60 * 5)
		for {
			FileListUpdate()
			time.Sleep(time.Second * 60 * 5)
		}
	}()

}

func FileListUpdate(c ...*pterm.SpinnerPrinter) {
	var err error
	var msg string
	defer func() {

		if c == nil {
			if msg != "" {
				fmt.Println(msg)
			}
			return
		}
		if err != nil {
			c[0].Fail(err)
		}
		c[0].Success()
	}()

	files, err := ioutil.ReadDir(Setting.TargetDir)
	if err != nil {
		return
	}

	tmp := make(FileIndex)
	for _, file := range files {
		if sid, err := strconv.Atoi(strings.Replace(file.Name(), ".osz", "", -1)); err == nil {
			tmp[sid] = file.ModTime()
		}
	}
	FileList = tmp
	msg = fmt.Sprintf(
		"%s File List Indexing : %s files\n",
		time.Now().Format("2006-01-02 15:04:05"),
		[10]string{strconv.Itoa(len(FileList))},
	)

}
