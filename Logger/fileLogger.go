package Logger

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/pterm/pterm"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var file *os.File
var Ch = make(chan struct{}) //UpdateLogFile

const (
	logPath        = "./log"
	maxLogFileSize = int64(1024 * 1024 * 1024)
)

func init() {
	go func() {
		setLogFile()
		checkLogFileLimit()
		_ = gocron.Every(1).Days().At("00:00:00").Do(setLogFile)
		_ = gocron.Every(1).Hours().At("00:00").Do(checkLogFileLimit)
		<-gocron.Start()
		//for {
		//	setLogFile()
		//	checkLogFileLimit()
		//	st, _ := time.Parse("20060102",time.Now().Add(time.Hour*24).Format("20060102"))
		//	st = st.Add(-time.Hour *9)
		//	time.Sleep(time.Duration(st.Unix() - time.Now().UTC().Unix())*time.Second)
		//}
	}()
	pterm.Info.Println("logfile Scheduler Started.")

}

func checkLogFileLimit() {
	checkDir()

	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		pterm.Error.Println(err)
		return
	}

	sort.Slice(files, func(i, j int) (tf bool) {
		fii, _ := strconv.Atoi(strings.Split(files[i].Name(), ".")[0])
		fij, _ := strconv.Atoi(strings.Split(files[j].Name(), ".")[0])
		return fii > fij
	})
	var fileSize int64
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ri, err := regexp.Match("^([0-9][.]log)$", []byte(f.Name()))
		if err != nil || ri {
			continue
		}
		fileSize += f.Size()

		if maxLogFileSize < fileSize {
			err := os.Remove(fmt.Sprintf("%s/%s", logPath, f.Name()))
			if err != nil {
				pterm.Error.Println(err)

			} else {
				pterm.Info.Printf("logfile %s Deleted.", f.Name())
			}

		}
	}
}

func checkDir() {
	if _, e := os.Stat(logPath); os.IsNotExist(e) {
		err := os.MkdirAll(logPath, 666)
		if err != nil {
			pterm.Error.Println(err)
			panic(err)
		}
	}
}

func setLogFile() {

	if file != nil {
		file.Close()
	}
	checkDir()

	fileName := fmt.Sprintf("%s/%s.log", logPath, time.Now().Format("060102"))
	pterm.Info.Println("SET LOG FILE: ", fileName)
	fpLog, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_SYNC, 0777)
	if err != nil {
		pterm.Error.Println(err)
	}
	file = fpLog

	log.SetOutput(file)
	Ch <- struct{}{}

}
