package Logger

import (
	"io"
	"os"
)

var file *os.File

func LoadLogger(writer *io.Writer) (){
	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	//Logger = log.New(fpLog)
}