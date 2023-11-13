package mtgapi

import (
	"github.com/gravestench/runtime"

	"github.com/gravestench/mtg/pkg/services/configFile"
)

type Dependency = MTGApiClient

type MTGApiClient interface {
	runtime.Service
	runtime.HasLogger
	runtime.HasDependencies
	configFile.HasDefaultConfig
	GetTypes() ([]string, error)
	GetSubTypes() ([]string, error)
	GetFormats() ([]string, error)
	GetSuperTypes() ([]string, error)
}
