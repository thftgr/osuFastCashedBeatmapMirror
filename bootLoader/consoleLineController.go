package bootLoader

import (
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
)

func BootMirror() {




	spinner, _ :=  pterm.DefaultSpinner.Start("Load Config File.")





	//ch := make(chan struct{})
	src.LoadConfig(spinner)
	spinner, _ =  pterm.DefaultSpinner.Start("Load Beatmap Files.")
	src.StartIndex(spinner)
	//spinner, _ =  pterm.DefaultSpinner.Start("[ waiting ] [ 0% ] Load RDBMS.")
	//db.ConnectMaria()
	//spinner, _ =  pterm.DefaultSpinner.Start("[ waiting ] [ 0% ] Load Bancho Token.")
	//spinner, _ =  pterm.DefaultSpinner.Start("[ waiting ] [ 0% ] Load Bancho Data API.")
	//go src.LoadBancho(ch)
	//db.ConnectMaria()
	//go Logger.LoadLogger(&LogIO)
	//_ = <-ch
	//src.RunGetBeatmapDataASBancho()
}
