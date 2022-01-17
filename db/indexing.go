package db

import (
	"github.com/pterm/pterm"
	"strconv"
	"strings"
	"time"
)

var Index = map[string][]int{}

type searchIndexTDS struct {
	id   int
	data string
}

func LoadIndex() {
	for {
		pterm.Info.Println("started database indexing")
		rows, err := Maria.Query(`select beatmapset_id, concat_ws(' ',artist, creator, title, tags) from osu.beatmapset order by beatmapset_id desc;`)
		if err != nil {
			pterm.Error.Println(err)
		}
		defer rows.Close()
		var tds searchIndexTDS
		for rows.Next() {
			err = rows.Scan(&tds.id, &tds.data)
			if err != nil {
				pterm.Error.Println(err)
				tds = searchIndexTDS{}
				return
			}
			data := strings.Split(strings.ToLower(strings.TrimSpace(tds.data)), " ")
			dataSize := len(data)
			Index[strconv.Itoa(tds.id)] = append(Index[strconv.Itoa(tds.id)], tds.id)
			for i := 0; i < dataSize; i++ {
				if data[i] != "" {
					Index[data[i]] = append(Index[data[i]], tds.id)
				}
			}
		}
		for key, val := range Index {
			Index[key] = makeSliceUnique(val)
		}

		pterm.Success.Println("end database indexing", len(Index))
		time.Sleep(time.Minute * 10)
	}
}

func SearchIndex(q string) (d []int) {
	data := strings.Split(strings.ToLower(strings.TrimSpace(q)), " ")
	dataSize := len(data)
	var ids = map[int]int{}
	for i := 0; i < dataSize; i++ {
		//fmt.Println(data[i], Index[data[i]])
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
	return
}

func makeSliceUnique(s []int) []int {
	keys := make(map[int]struct{})
	res := make([]int, 0)
	for _, val := range s {
		if _, ok := keys[val]; ok {
			continue
		} else {
			keys[val] = struct{}{}
			res = append(res, val)
		}
	}
	return res
}
