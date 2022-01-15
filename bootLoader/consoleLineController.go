package bootLoader

import (
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/pterm/pterm"
	"os"
)

func init() {
	ch := make(chan struct{})

	spinner, _ := pterm.DefaultSpinner.Start("Load Config File.")
	config.LoadConfig(spinner)

	spinner, _ = pterm.DefaultSpinner.Start("Load Beatmap Files.")
	src.StartIndex(spinner)

	spinner, _ = pterm.DefaultSpinner.Start("Load RDBMS.")
	db.ConnectMaria()

	go src.LoadBancho(ch)
	_ = <-ch
	if os.Getenv("debug") != "true" {
		go src.RunGetBeatmapDataASBancho()
	}
}
