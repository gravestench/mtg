package scryfall

import (
	"image"

	"github.com/BlueMonday/go-scryfall"
	"github.com/gravestench/runtime"

	"github.com/gravestench/mtg/pkg/services/configFile"
)

type Dependency = ScryfallClient

type ScryfallClient interface {
	runtime.Service
	runtime.HasLogger
	runtime.HasDependencies
	configFile.HasDefaultConfig
	Search(name string) (*scryfall.CardListResponse, error)
	SearchWithDeckList(list string) []scryfall.Card
	GetImagesFromCard(card scryfall.Card) ([]image.Image, error)
	GetImagesFromDeckList(list string) ([]image.Image, error)
}
