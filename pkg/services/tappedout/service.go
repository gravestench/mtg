package tappedout

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BlueMonday/go-scryfall"
	"github.com/gravestench/runtime"
	"github.com/rs/zerolog"

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
	return "TappedOut"
}

func (s *Service) BindLogger(logger *zerolog.Logger) {
	s.logger = logger
}

func (s *Service) Logger() *zerolog.Logger {
	return s.logger
}

func (s *Service) ConfigFileName() string {
	return "tappedout.json"
}

func (s *Service) DefaultConfig() (cfg configFile.Config) {
	cfg.Group("tappedout").Set("directory", "/tmp")

	return
}

func (s *Service) GetDeckList(uri string) (string, error) {
	const decksEndpoint = "https://tappedout.net/mtg-decks"

	uri = strings.Trim(uri, "/ ")
	uri = strings.ReplaceAll(uri, "http://", "https://")

	if !strings.HasPrefix(uri, decksEndpoint) {
		uri = strings.Join([]string{decksEndpoint, uri}, "/")
	}

	xmlBytes, err := getXML(uri)
	if err != nil {
		return "", fmt.Errorf("failed to get XML: %v", err)
	}

	const regexGrabDeckList = `<textarea id="mtga-textarea">(([^<])+\n?\t*)+`
	r := regexp.MustCompile(regexGrabDeckList)

	matches := r.FindStringSubmatch(string(xmlBytes))

	if len(matches) < 2 {
		return "", fmt.Errorf("could not find deck list on page")
	}

	listBytes := matches[1]

	dstDir := s.cfg.Group("tappedout").GetString("directory")
	uriParts := strings.Split(uri, "/")
	fileName := uriParts[len(uriParts)-1]
	decklistPath := filepath.Join(dstDir, fmt.Sprintf("%s.txt", fileName))

	if _, err = os.Stat(dstDir); os.IsNotExist(err) {
		err = os.Mkdir(dstDir, 0755)
		// TODO: handle error
	}

	err = os.WriteFile(decklistPath, []byte(listBytes), 0777)
	if err != nil {
		return "", fmt.Errorf("writing deck list from tappedout.net: %v", err)
	}

	return listBytes, nil
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
