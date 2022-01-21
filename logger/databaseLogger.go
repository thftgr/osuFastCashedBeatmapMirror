package logger

import (
	"encoding/json"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/pterm/pterm"
	"time"
)

//var LogBuffer bytes.Buffer

//func LoadLogger() {
//	time.Sleep(time.Second)
//	for {
//		if b.Len() < 1 {
//			time.Sleep(time.Second)
//			continue
//		}
//		line, err := bufio.NewReader(b).ReadBytes(0x0A)
//		if err != nil {
//			pterm.Error.Println(err)
//			continue
//		}
//		js := map[string]interface{}{}
//		if err = json.Unmarshal(line, &js); err != nil {
//			pterm.Error.Println(err)
//			continue
//		}
//		//fmt.Println(string(line))
//		t, err := time.Parse(time.RFC3339Nano, js["time"].(string))
//		if err != nil {
//			pterm.Error.Println(err)
//			continue
//		}
//
//		_, err = db.Maria.Exec(
//			`INSERT INTO BeatmapMirror.api_log (time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?);`,
//			t.Format("2006-01-02 15-04-05"),
//			js["id"],
//			js["remote_ip"],
//			js["host"],
//			js["method"],
//			js["uri"],
//			js["user_agent"],
//			js["status"],
//			js["error"],
//			js["latency"],
//			js["latency_human"],
//			js["bytes_in"],
//			js["bytes_out"],
//		)
//		//time, request_id, remote_ip, host, method, uri, user_agent, status, error, latency, latency_human, bytes_in, bytes_out
//
//		if err != nil {
//			pterm.Error.Println(err)
//			continue
//		}
//	}
//
//}

func Error(v *bodyStruct.ErrorStruct) (vv *bodyStruct.ErrorStruct) {
	go func() {
		b, _ := json.Marshal(v)
		pterm.Error.Println(time.Now().Format("2006-01-02 15:04:05"), string(b))
	}()

	//TODO DB 에 저장
	return v

}
func Info(v *bodyStruct.ErrorStruct) (vv *bodyStruct.ErrorStruct) {
	go func() {
		b, _ := json.Marshal(v)
		pterm.Info.Println(time.Now().Format("2006-01-02 15:04:05"), string(b))
	}()

	//TODO DB 에 저장
	return v

}
