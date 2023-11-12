package tappedout

import (
	"github.com/gravestench/runtime"

	"github.com/gravestench/mtg/pkg/services/configFile"
)

type Dependency = ScryfallClient

type ScryfallClient interface {
	runtime.Service
	runtime.HasLogger
	runtime.HasDependencies
	configFile.HasDefaultConfig
	GetDeckList(uri string) (string, error)
}
