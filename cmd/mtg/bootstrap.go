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

	camera := s.renderer.GetDefaultCamera()
	camera.Zoom = 0.25

	grid := 5
	pad := 0.05

	var ox, oy int
	for idx, img := range images {
		ox = img.Bounds().Dx()
		oy = img.Bounds().Dy() / 2
		padX, padY := float64(img.Bounds().Dx())*pad, float64(img.Bounds().Dy())*pad
		gx := idx % grid
		gy := idx / grid
		x, y := int(padX)+ox+gx*img.Bounds().Dx(), int(padY)+oy+gy*img.Bounds().Dy()
		node := s.renderer.NewRenderable()
		node.SetImage(img)
		node.SetPosition(float32(x), float32(y))
		node.SetRotation(float32(rand.Intn(6000)-3000) / 1000)
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
