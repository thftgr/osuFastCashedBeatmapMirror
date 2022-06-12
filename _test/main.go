package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/goccy/go-json"
	"golang.org/x/net/html"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var beatmapSetUrlRegex, _ = regexp.Compile(`(?:https:\/\/osu[.]ppy[.]sh\/beatmapsets\/)([0-9].+?)$`)
var beatmapPacks []packType

func main() {

	doc, err := htmlquery.LoadURL("https://osu.ppy.sh/beatmaps/packs")
	if err != nil {
		panic(err)
	}
	beatmapPacks = getTypes(doc)
	for i := 0; i < len(beatmapPacks); i++ {
		beatmapPacks[i].setMaxPage()
	}

	for i := 0; i < len(beatmapPacks); i++ {
		for j := 0; j <= beatmapPacks[i].MaxPage; j++ { //https://osu.ppy.sh/beatmaps/packs?type=standard&page=2
			url := fmt.Sprintf("%s&page=%d", beatmapPacks[i].Url, j)
			log.Println(url)
			doc, err := htmlquery.LoadURL(url)
			if err != nil {
				panic(err)
			}
			packd := getPackData(doc)
			for k := 0; k < len(packd); k++ {
				packd[k].setAllBeatmaps()
				time.Sleep(time.Second)
			}
			beatmapPacks[i].Packs = append(beatmapPacks[i].Packs, packd...)
			time.Sleep(time.Second)
		}

	}

	log.Println(ToJsonString(beatmapPacks))
	f, err := os.OpenFile("./packs.json", os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_SYNC, 0777)
	if err == nil {
		_, _ = f.WriteString(ToJsonString(beatmapPacks))
		_ = f.Close()
	} else {
		log.Fatal(err)
	}

}

type pack struct {
	Url      string `json:"url"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	By       string `json:"by"`
	Beatmaps []int
}

type packType struct {
	Type    string
	Url     string
	MaxPage int
	Packs   []pack
}

func (v *pack) setAllBeatmaps() {
	url := v.Url + "/raw"
	log.Println(url)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		panic(err)
	}
	list, _ := htmlquery.QueryAll(doc, `//a[@class="beatmap-pack-items__link"]`)

	for _, n := range list {
		rawUrl := htmlquery.SelectAttr(n, "href")
		if !beatmapSetUrlRegex.MatchString(rawUrl) {
			continue
		}
		str := beatmapSetUrlRegex.FindAllStringSubmatch(rawUrl, 1)
		p, err := strconv.Atoi(str[0][1])
		if err != nil {
			v.Beatmaps = append(v.Beatmaps, p)
		}

	}
}

func (v *packType) setMaxPage() {
	doc, err := htmlquery.LoadURL(v.Url)
	if err != nil {
		panic(err)
	}
	list, _ := htmlquery.QueryAll(doc, `//a[@class="pagination-v2__link"]`)

	for _, n := range list {
		p, _ := strconv.Atoi(htmlquery.InnerText(n))
		if p > v.MaxPage {
			v.MaxPage = p
		}
	}
}

func getTypes(node *html.Node) (packs []packType) {
	list, _ := htmlquery.QueryAll(node, `//a[contains(@class,"page-mode-link")]`)
	for _, n := range list {
		packs = append(packs, packType{
			Url: htmlquery.SelectAttr(n, "href"),
		})
	}
	return
}

func getPackData(node *html.Node) (packs []pack) {
	list, _ := htmlquery.QueryAll(node, `//div[@class="beatmap-pack js-beatmap-pack js-accordion__item"]`)
	for _, n := range list {
		url, _ := htmlquery.Query(n, `.//a[@class="beatmap-pack__header js-accordion__item-header"]`)
		name, _ := htmlquery.Query(n, `.//div[@class="beatmap-pack__name"]`)
		id, _ := htmlquery.Query(n, `.//div[@class="beatmap-pack js-beatmap-pack js-accordion__item"]`)
		date, _ := htmlquery.Query(n, `.//span[@class="beatmap-pack__date"]`)
		by, _ := htmlquery.Query(n, `.//span[@class="beatmap-pack__author beatmap-pack__author--bold"]`)

		packs = append(packs, pack{
			Url:  htmlquery.SelectAttr(url, "href"),
			Name: htmlquery.InnerText(name),
			Id:   htmlquery.SelectAttr(id, "data-pack-id"),
			Date: htmlquery.InnerText(date),
			By:   htmlquery.InnerText(by),
		})
	}
	return
}

func ToJsonString(i interface{}) (str string) {
	b, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	return string(b)
}
