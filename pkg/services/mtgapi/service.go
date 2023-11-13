package mtgapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/BlueMonday/go-scryfall"
	"github.com/gravestench/runtime"
	"github.com/rs/zerolog"

	"github.com/MagicTheGathering/mtg-sdk-go"

	"github.com/gravestench/mtg/pkg/services/configFile"
)

type Service struct {
	client     *scryfall.Client
	logger     *zerolog.Logger
	cfgManager configFile.Dependency
	cfg        *configFile.Config
}

func (s *Service) DependenciesResolved() bool {
	if s.cfgManager == nil {
		return false
	}

	return true
}

func (s *Service) ResolveDependencies(r runtime.R) {
	for _, service := range r.Services() {
		if candidate, ok := service.(configFile.Dependency); ok {
			s.cfgManager = candidate
		}
	}
}

func (s *Service) Init(rt runtime.Runtime) {
	client, err := scryfall.NewClient()
	if err != nil {
		s.logger.Fatal().Msgf("could not open scryfall client: %v", err)
	}

	s.client = client

	cfg, err := s.cfgManager.GetConfigByFileName(s.ConfigFileName())
	if err != nil {
		s.logger.Fatal().Msgf("loading config file: %v", err)
	}

	s.cfg = cfg
}

func (s *Service) Name() string {
	return "Scryfall"
}

func (s *Service) BindLogger(logger *zerolog.Logger) {
	s.logger = logger
}

func (s *Service) Logger() *zerolog.Logger {
	return s.logger
}

func (s *Service) ConfigFileName() string {
	return "scryfall.json"
}

func (s *Service) DefaultConfig() (cfg configFile.Config) {
	cfg.Group("scryfall").Set("directory", "/tmp")

	return
}

func (s *Service) GetTypes() ([]string, error) {
	return mtg.GetTypes()
}

func (s *Service) GetSubTypes() ([]string, error) {
	return mtg.GetSubTypes()
}

func (s *Service) GetFormats() ([]string, error) {
	return mtg.GetFormats()
}

func (s *Service) GetSuperTypes() ([]string, error) {
	return mtg.GetSuperTypes()
}

func download(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("issuing http request: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	d, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return d, err
}
