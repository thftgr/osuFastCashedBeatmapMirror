package db

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/timeUnit"
	"github.com/pterm/pterm"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	//Index       = map[string][]int{}
	SearchCache = map[string]map[byte][]int{}
	IdInt       = map[int][]int{}
	artist      = map[string][]int{}
	creator     = map[string][]int{}
	title       = map[string][]int{}
	tags        = map[string][]int{}
	mutex       = sync.RWMutex{}
)

type searchIndexTDS struct {
	id      int
	artist  string
	creator string
	title   string
	tags    string
}

func LoadIndex() {
	for {
		doIndex()
		debug.FreeOSMemory()
		time.Sleep(time.Minute * 10)
	}
}

//TODO 맵셋 인덱싱 + 맵 인덱싱
//TODO 기존 3433ms
func doIndex() {
	var countKey = 0
	var countValue = 0
	pterm.Info.Println(timeUnit.GetTime(), "started database indexing")
	st := time.Now().UnixMilli()

	//select beatmap_id, beatmapset_id from osu.beatmap;
	rows, err := Maria.Query(`select beatmap_id, beatmapset_id from osu.beatmap;`)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var m, s int
		err = rows.Scan(&m, &s)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		mutex.Lock()
		IdInt[m] = append(IdInt[m], s)
		mutex.Unlock()
	}

	rows, err = Maria.Query(`select beatmapset_id, artist, creator, title, tags from osu.beatmapset`)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tds := searchIndexTDS{}
		err = rows.Scan(&tds.id, &tds.artist, &tds.creator, &tds.title, &tds.tags)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		mutex.Lock()
		IdInt[tds.id] = append(IdInt[tds.id], tds.id)
		mutex.Unlock()
		// 타입별 쓰레드 ============================================================
		var ch = make(chan struct{})
		go func() {
			data := strings.Split(strings.ToLower(strings.TrimSpace(tds.artist)), " ")
			dataSize := len(data)
			for i := 0; i < dataSize; i++ {
				if d := data[i]; d != "" {
					mutex.Lock()
					artist[d] = append(artist[d], tds.id)
					mutex.Unlock()
				}
			}
			ch <- struct{}{}
		}()
		go func() {
			data := strings.Split(strings.ToLower(strings.TrimSpace(tds.creator)), " ")
			dataSize := len(data)
			for i := 0; i < dataSize; i++ {
				if d := data[i]; d != "" {
					mutex.Lock()
					creator[d] = append(creator[d], tds.id)
					mutex.Unlock()
				}
			}
			ch <- struct{}{}
		}()
		go func() {
			data := strings.Split(strings.ToLower(strings.TrimSpace(tds.title)), " ")
			dataSize := len(data)
			for i := 0; i < dataSize; i++ {
				if d := data[i]; d != "" {
					mutex.Lock()
					title[d] = append(title[d], tds.id)
					mutex.Unlock()
				}
			}
			ch <- struct{}{}
		}()
		go func() {
			data := strings.Split(strings.ToLower(strings.TrimSpace(tds.tags)), " ")
			dataSize := len(data)
			for i := 0; i < dataSize; i++ {
				if d := data[i]; d != "" {
					mutex.Lock()
					tags[d] = append(tags[d], tds.id)
					mutex.Unlock()
				}
			}
			ch <- struct{}{}
		}()

		for i := 0; i < 4; i++ {
			_ = <-ch
		}
	}
	// 중복제거 쓰레드 ============================================================
	var ch2 = make(chan struct{})
	go func() {
		for key, val := range IdInt {
			mutex.Lock()
			IdInt[key] = *makeSliceUnique(val)
			mutex.Unlock()
			countValue += len(IdInt[key])
			countKey += int(unsafe.Sizeof([]int{})) * len(IdInt[key])
		}
		ch2 <- struct{}{}
	}()
	go func() {
		for key, val := range artist {
			mutex.Lock()
			artist[key] = *makeSliceUnique(val)
			mutex.Unlock()
			countValue += len(artist[key])
			countKey += len([]byte(key))
		}
		ch2 <- struct{}{}
	}()
	go func() {
		for key, val := range creator {
			mutex.Lock()
			creator[key] = *makeSliceUnique(val)
			mutex.Unlock()
			countValue += len(creator[key])
			countKey += len([]byte(key))
		}
		ch2 <- struct{}{}
	}()
	go func() {
		for key, val := range title {
			mutex.Lock()
			title[key] = *makeSliceUnique(val)
			mutex.Unlock()
			countValue += len(title[key])
			countKey += len([]byte(key))
		}
		ch2 <- struct{}{}
	}()
	go func() {
		for key, val := range tags {
			mutex.Lock()
			tags[key] = *makeSliceUnique(val)
			mutex.Unlock()
			countValue += len(tags[key])
			countKey += len([]byte(key))
		}
		ch2 <- struct{}{}
	}()

	for i := 0; i < 5; i++ {
		_ = <-ch2
	}

	//============================================================
	et := time.Now().UnixMilli()
	a := int(unsafe.Sizeof(map[string][]int{})) * countKey
	b := int(unsafe.Sizeof([]int{})) * countValue

	mutex.Lock()
	SearchCache = map[string]map[byte][]int{}
	mutex.Unlock()
	pterm.Success.Printfln("%s end database indexing %d keys. %d links.using %s of memory. %dms", timeUnit.GetTime(), countKey, countValue, calcMeoerySize(a+b), et-st)
}

func SearchIndex(q string, option byte) (d []int) {
	st := time.Now().UnixMilli()
	t := SearchCache[q][option]
	if len(t) > 0 {
		return t
	}
	data := strings.Split(strings.ToLower(strings.TrimSpace(q)), " ")
	dataSize := len(data)
	var ids = map[int]int{}
	var allIds []int
	for i := 0; i < dataSize; i++ {
		di := data[i]
		if di != "" {
			mutex.Lock()
			if option&0x01 > 0 {

				allIds = append(allIds, artist[di]...)
			}
			if option&0x02 > 0 {
				allIds = append(allIds, creator[di]...)
			}
			if option&0x04 > 0 {
				allIds = append(allIds, tags[di]...)
			}
			if option&0x08 > 0 {
				allIds = append(allIds, title[di]...)
			}
			if option&0x16 > 0 {
				dii, err := strconv.Atoi(di)
				if err != nil {

					allIds = append(allIds, IdInt[dii]...)

				}
			}
			mutex.Unlock()

		}
	}
	if allIds != nil {
		ds := len(allIds)
		tv := allIds
		for j := 0; j < ds; j++ {
			ids[tv[j]] += 1
		}
	}
	for k, v := range ids {
		if v == dataSize {
			d = append(d, k)
		}
	}
	mutex.Lock()
	SearchCache[q] = map[byte][]int{option: d}
	mutex.Unlock()
	et := time.Now().UnixMilli()
	fmt.Println(et-st, "ms")
	return
}

func makeSliceUnique(s []int) *[]int {
	keys := make(map[int]struct{})
	res := make([]int, 0)
	l := len(s)
	for i := 0; i < l; i++ {
		keys[s[i]] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	keys = nil
	return &res
}

//func makeSliceUnique(s *[]int) *[]int {
//	keys := make(map[int]struct{})
//	res := make([]int, 0)
//	for _, val := range *s {
//		if _, ok := keys[val]; ok {
//			continue
//		} else {
//			keys[val] = struct{}{}
//			res = append(res, val)
//		}
//	}
//	keys = nil
//	return &res
//}
func calcMeoerySize(l int) (s string) {
	size := float32(l)
	if runtime.GOOS == "windows" {
		if size > 1099511627776 { //TB
			return fmt.Sprintf("%.3f%s", size/1099511627776, "TB")
		} else if size > 1073741824 { //GB
			return fmt.Sprintf("%.3f%s", size/1073741824, "GB")
		} else if size > 1048576 { //MB
			return fmt.Sprintf("%.3f%s", size/1048576, "MB")
		} else if size > 1024 { //KB
			return fmt.Sprintf("%.3f%s", size/1024, "KB")
		}
	} else {
		if size > 1000000000000 { //TB
			return fmt.Sprintf("%.3f%s", size/1000000000000, "TB")
		} else if size > 1000000000 { //GB
			return fmt.Sprintf("%.3f%s", size/1000000000, "GB")
		} else if size > 1000000 { //MB
			return fmt.Sprintf("%.3f%s", size/1000000, "MB")
		} else if size > 1000 { //KB
			return fmt.Sprintf("%.3f%s", size/1000, "KB")
		}
	}

	return fmt.Sprintf("%.3f%s", size, "B")
}
