package main

import (
	"encoding/json"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Point struct {
	X int
	Y int
}

func (v *Point) Equal(p Point) bool {
	return v.X == p.X && v.Y == p.Y
}
func main() {

	MAP := [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 1, 0, 1, 0, 1},
		{1, 1, 0, 0, 0, 1, 0, 1},
		{1, 1, 1, 0, 1, 1, 0, 1},
		{1, 1, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 1, 1, 1},
		{1, 1, 0, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
	}
	START := Point{
		X: 1,
		Y: 1,
	}
	END := Point{
		X: 6,
		Y: 7,
	}
	line := Line{
		Map:   MAP,
		Start: START,
		End:   END,
		Pass:  nil,
	}
	FindPass(line)

}

type Line struct {
	Map   [][]int
	Start Point
	End   Point
	Pass  []Point
}

func FindPass(v Line) {
	org := v
	check := func(p Point) {
		if v.End.Equal(p) {
			log.Println(p.X, p.Y)
			panic("END")
		}
		if !p.Equal(v.Start) && v.Map[p.Y][p.X] == 0 {
			log.Println(p.X, p.Y, "|", v.Map[p.Y][p.X])
			v.Start = p
			v.Pass = append(v.Pass, p)
			log.Println(ToJsonString(v))
			FindPass(v)
		}
	}
	if v.Start.X > 0 {
		p := Point{
			X: v.Start.X - 1,
			Y: v.Start.Y,
		}
		check(p)

	}
	v = org
	if v.Start.X < len(v.Map)-1 {
		p := Point{
			X: v.Start.X + 1,
			Y: v.Start.Y,
		}
		check(p)
	}
	v = org
	if v.Start.Y > 0 {
		p := Point{
			X: v.Start.X,
			Y: v.Start.Y - 1,
		}
		check(p)
	}
	v = org
	if v.Start.Y < len(v.Map[v.Start.X])-1 {
		p := Point{
			X: v.Start.X,
			Y: v.Start.Y + 1,
		}
		check(p)
	}

}
func ToJsonString(i interface{}) (str string) {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
		return
	}
	return string(b)
}
