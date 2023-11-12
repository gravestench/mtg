package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gravestench/runtime"
	"github.com/rs/zerolog"

	"github.com/gravestench/mtg/pkg/services/raylibRenderer"
	"github.com/gravestench/mtg/pkg/services/scryfall"
	"github.com/gravestench/mtg/pkg/services/tappedout"
)

type bootstrap struct {
	logger    *zerolog.Logger
	renderer  raylibRenderer.Dependency
	scryfall  scryfall.Dependency
	tappedout tappedout.Dependency
}

func (s *bootstrap) BindLogger(logger *zerolog.Logger) {
	s.logger = logger
}

func (s *bootstrap) Logger() *zerolog.Logger {
	return s.logger
}

func (s *bootstrap) DependenciesResolved() bool {
	if s.renderer == nil {
		return false
	}

	if s.scryfall == nil {
		return false
	}

	return true
}

func (s *bootstrap) ResolveDependencies(r runtime.R) {
	for _, service := range r.Services() {
		switch candidate := service.(type) {
		case raylibRenderer.Dependency:
			s.renderer = candidate
		case scryfall.Dependency:
			s.scryfall = candidate
		case tappedout.Dependency:
			s.tappedout = candidate
		}
	}
}

func (s *bootstrap) Init(r runtime.R) {
	const deckURL = "https://tappedout.net/mtg-decks/a-slow-painful-death/"

	list, err := s.tappedout.GetDeckList(deckURL)
	if err != nil {
		s.logger.Error().Msgf("getting deck list from tappedout: %v", err)
		return
	}

	images, err := s.scryfall.GetImagesFromDeckList(list)
	if len(images) < 1 {
		return
	}

	// zoom out instead of adjusting all of the card dimensions
	camera := s.renderer.GetDefaultCamera()
	camera.Zoom = 0.25

	const (
		gridWidth             = 5
		gridCellPadNormalized = 0.05
	)

	var gridOriginX, gridOriginY int
	for idx, img := range images {
		// offset of grid, we only need to set this once,
		// but all the card dimensions are the same
		gridOriginX, gridOriginY = img.Bounds().Dx(), img.Bounds().Dy()/2

		// padding between the cards, relative to the card dimensions
		padX := float64(img.Bounds().Dx()) * gridCellPadNormalized
		padY := float64(img.Bounds().Dy()) * gridCellPadNormalized

		// grid cell x,y index, based off of card index
		cellX, cellY := idx%gridWidth, idx/gridWidth

		// actual x,y position of the card, based off of the grid origin,
		// padding, and dimensions of the card
		x, y := int(padX)+gridOriginX+cellX*img.Bounds().Dx(), int(padY)+gridOriginY+cellY*img.Bounds().Dy()

		// create our renderable object and set its properties
		node := s.renderer.NewRenderable()
		node.SetImage(img)
		node.SetPosition(float32(x), float32(y))

		// animate the rotation, with a random starting phase for each card,
		// and set up an OnUpdate callback which applies new rotation based on
		// the current timestamp
		phase := (rand.Float64() * 2) - 1
		node.OnUpdate(func() {
			rotation := float32(math.Sin((float64(time.Now().UnixNano()) / 250000000.0) + phase))
			node.SetRotation(rotation)
		})
	}
}

func (s *bootstrap) Name() string {
	return "Bootstrap"
}
