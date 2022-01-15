package bootLoader

import (
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/config"
	"github.com/thftgr/osuFastCashedBeatmapMirror/db"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
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
