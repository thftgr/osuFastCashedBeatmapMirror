package Logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"time"
)



func LoadLogger(b * bytes.Buffer) {
	time.Sleep(time.Second)
	for {
		if b.Len() < 1 {
			time.Sleep(time.Second)
			continue
		}
		line, err := bufio.NewReader(b).ReadBytes(0x0A)
		if err != nil {
			fmt.Println( err)
			continue
		}
		js := map[string]interface{}{}
		if err = json.Unmarshal(line, &js); err != nil {
			fmt.Println( err)
			continue
		}
		//fmt.Println(string(line))
		t, err := time.Parse(time.RFC3339Nano, js["time"].(string))
		if err != nil {
			fmt.Println( err)
			continue
		}
		//{"time":"2021-07-02T22:11:37.391019+09:00","id":"De14R","remote_ip":"::1","host":"localhost","method":"GET","uri":"/search","user_agent":"","status":200,"error":"","latency":489638000,"latency_human":"489.638ms","bytes_in":0,"bytes_out":218561}

		//time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out
		err = src.InsertAPILog(
			t.Format("2006-01-02 15-04-05"),
			js["id"],
			js["remote_ip"],
			js["host"],
			js["method"],
			js["uri"],
			js["user_agent"],
			js["status"],
			js["error"],
			js["latency"],
			js["latency_human"],
			js["bytes_in"],
			js["bytes_out"],
		)
		if err != nil {
			fmt.Println( err)
			continue
		}

	}

}