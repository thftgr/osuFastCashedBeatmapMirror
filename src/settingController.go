package src

import (
	"encoding/json"
	"github.com/pterm/pterm"
	"io/ioutil"
)

type config struct {
	Port      string `json:"port"`
	TargetDir string `json:"targetDir"`
	Sql       struct {
		Id     string `json:"id"`
		Passwd string `json:"passwd"`
		Url    string `json:"url"`
	} `json:"sql"`
	Osu struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
		Token    struct {
			TokenType    string `json:"token_type"`
			ExpiresIn    int64  `json:"expires_in"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"token"`
		BeatmapUpdate struct {
			UpdatedAsc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_asc"`
			UpdatedDesc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"updated_desc"`
			GraveyardAsc struct {
				LastUpdate string `json:"last_update"`
				Id         string `json:"_id"`
			} `json:"graveyard_asc"`
		} `json:"beatmapUpdate"`
	} `json:"osu"`
}

var Setting config

func LoadConfig(c *pterm.SpinnerPrinter) {
	var err error
	defer func() {
		if err != nil {
			c.Fail(err)
		}
		c.Success()
	}()
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &Setting)
	if err != nil {
		return
	}

}
func (v *config) Save() {
	file, _ := json.MarshalIndent(v, "", "  ")
	_ = ioutil.WriteFile("config.json", file, 0755)
}
