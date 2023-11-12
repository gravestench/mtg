package main

import (
	"github.com/faiface/mainthread"
	"github.com/gravestench/runtime"

	"github.com/gravestench/mtg/pkg/services/cacheManager"
	"github.com/gravestench/mtg/pkg/services/configFile"
	"github.com/gravestench/mtg/pkg/services/raylibRenderer"
	"github.com/gravestench/mtg/pkg/services/scryfall"
	"github.com/gravestench/mtg/pkg/services/tappedout"
)

func main() {
	rt := runtime.New("MTG")

	rt.Add(&cacheManager.Service{})
	rt.Add(&configFile.Service{RootDirectory: "~/.config/mtg"})
	rt.Add(&raylibRenderer.Service{})
	rt.Add(&scryfall.Service{})
	rt.Add(&tappedout.Service{})
	rt.Add(&bootstrap{})

	mainthread.Run(rt.Run)
}
