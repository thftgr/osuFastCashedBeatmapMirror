package db

import (
	"github.com/pterm/pterm"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var Index = map[string][]int{}
var SearchCache = map[string][]int{}

type searchIndexTDS struct {
	id   int
	data string
}

func LoadIndex() {
	for {
		doIndex()
		debug.FreeOSMemory()
		time.Sleep(time.Minute * 10)
	}
}

//TODO 맵셋 인덱싱 + 맵 인덱싱
func doIndex() {

	pterm.Info.Println("started database indexing")
	rows, err := Maria.Query(`select beatmapset_id, concat_ws(' ',artist, creator, title, tags ) from osu.beatmapset order by beatmapset_id desc;`)
	if err != nil {
		pterm.Error.Println(err)
	}
	defer rows.Close()
	var tds searchIndexTDS
	for rows.Next() {
		tds = searchIndexTDS{}
		err = rows.Scan(&tds.id, &tds.data)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		data := strings.Split(strings.ToLower(strings.TrimSpace(tds.data)), " ")
		dataSize := len(data)
		id := strconv.Itoa(tds.id)
		Index[id] = append(Index[id], tds.id)
		for i := 0; i < dataSize; i++ {
			if data[i] != "" {
				Index[data[i]] = append(Index[data[i]], tds.id)
			}
		}
	}
	var countKey = 0
	var countValue = 0
	for key, val := range Index {
		Index[key] = *makeSliceUnique(&val)
		countValue += len(Index[key])
		countKey += len([]byte(key))
	}
	a := int(unsafe.Sizeof(map[string][]int{})) * countKey
	b := int(unsafe.Sizeof([]int{})) * countValue
	SearchCache = map[string][]int{}
	pterm.Success.Printfln("end database indexing %d keys. %d links.using %d bytes of memory", countKey, countValue, a+b)
}
func SearchIndex(q string) (d []int) {
	t := SearchCache[q]
	if len(t) > 0 {
		return t
	}
	data := strings.Split(strings.ToLower(strings.TrimSpace(q)), " ")
	dataSize := len(data)
	var ids = map[int]int{}
	for i := 0; i < dataSize; i++ {
		if data[i] != "" {
			if Index[data[i]] != nil {
				ds := len(Index[data[i]])
				tv := Index[data[i]]
				for j := 0; j < ds; j++ {
					ids[tv[j]] += 1
				}
			}
		}
	}
	for k, v := range ids {
		if v == dataSize {
			d = append(d, k)
		}
	}
	SearchCache[q] = d
	return
}

func makeSliceUnique(s *[]int) *[]int {
	keys := make(map[int]struct{})
	res := make([]int, 0)
	for _, val := range *s {
		if _, ok := keys[val]; ok {
			continue
		} else {
			keys[val] = struct{}{}
			res = append(res, val)
		}
	}
	keys = nil
	return &res
}
