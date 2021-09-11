package bootLoader

import (
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"os"
)

func init() {
	ch := make(chan struct{})

	spinner, _ :=  pterm.DefaultSpinner.Start("Load Config File.")
	src.LoadConfig(spinner)

	spinner, _ =  pterm.DefaultSpinner.Start("Load Beatmap Files.")
	src.StartIndex(spinner)

	spinner, _ =  pterm.DefaultSpinner.Start("Load RDBMS.")
	src.ConnectMaria(spinner)

	go src.LoadBancho(ch)
	_ = <-ch
	if os.Getenv("debug") != "true" {
		go src.RunGetBeatmapDataASBancho()
	}
}
