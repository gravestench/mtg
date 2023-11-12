package scryfall

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/BlueMonday/go-scryfall"
	"github.com/gravestench/runtime"
	"github.com/rs/zerolog"

	"github.com/gravestench/mtg/pkg/services/configFile"
)

const (
	regexSplitMtgArenaLine = `(?P<Count>\d+) (?P<Name>[^(]+) \((?P<Set>[^)]+)\)(?P<CollectorNumber> \d+)?`
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

func (s *Service) Search(name string) (*scryfall.CardListResponse, error) {
	sco := scryfall.SearchCardsOptions{
		Unique:        scryfall.UniqueModePrints,
		Order:         scryfall.OrderSet,
		Dir:           scryfall.DirDesc,
		IncludeExtras: true,
	}

	result, err := s.client.SearchCards(context.Background(), name, sco)
	if err != nil {
		return nil, fmt.Errorf("could not search card: %v", err)
	}

	return &result, nil
}

func (s *Service) SearchWithDeckList(list string) (cards []scryfall.Card) {
	const (
		regexSplitMtgArenaLine = `(?P<Count>\d+) (?P<Name>[^(]+) \((?P<Set>[^)]+)\)(?P<CollectorNumber> \d+)?`
	)

	var lines []string

	lines = append(lines, strings.Split(list, "\n")...)

	s.logger.Info().Msgf("processing cards...")

	mtgArenaLineSplitter := regexp.MustCompile(regexSplitMtgArenaLine)

	for _, line := range lines {
		line = strings.Trim(line, "\n\t ")
		if line == "" {
			continue
		}

		match := mtgArenaLineSplitter.FindStringSubmatch(line)

		data := make(map[string]string)
		for i, name := range mtgArenaLineSplitter.SubexpNames() {
			if i > 0 && i <= len(match) {
				data[name] = match[i]
			}
		}

		name := data["Name"]
		name = strings.Split(name, " // ")[0]

		result, err := s.Search(name)
		if err != nil {
			s.logger.Error().Msgf("searching scryfall for %q: %v", name, err)
			continue
		}

		if len(result.Cards) < 1 {
			s.logger.Warn().Msgf("no cards found for `%v`", name)
			continue
		}

		cards = append(cards, result.Cards...)
	}

	return
}

func (s *Service) scryfallGetFirstMatchCardsFromDeckList(list string, cards []scryfall.Card) (result []scryfall.Card) {
	lines := strings.Split(list, "\n")

	mtgArenaLineSplitter := regexp.MustCompile(regexSplitMtgArenaLine)

	for _, line := range lines {
		line = strings.Trim(line, "\n\t ")
		if line == "" {
			continue
		}

		match := mtgArenaLineSplitter.FindStringSubmatch(line)

		deckListLineData := make(map[string]string)
		for i, key := range mtgArenaLineSplitter.SubexpNames() {
			if i > 0 && i <= len(match) {
				deckListLineData[key] = match[i]
			}
		}

		name := deckListLineData["Name"]
		name = strings.Split(name, " // ")[0]

		for _, card := range cards {
			if strings.ToLower(card.Name) != strings.ToLower(deckListLineData["Name"]) {
				continue
			}

			if strings.ToLower(card.Set) != strings.ToLower(deckListLineData["Set"]) {
				continue
			}

			//if deckListLineData["CollectorNumber"] != "" {
			//	if strings.ToLower(card.CollectorNumber) != strings.ToLower(deckListLineData["CollectorNumber"]) {
			//		continue
			//	}
			//}

			result = append(result, card)

			break
		}
	}

	return
}

func (s *Service) GetImagesFromCard(card scryfall.Card) (images []image.Image, err error) {
	collectorNumber := card.CollectorNumber
	if collectorNumber != "" {
		collectorNumber = fmt.Sprintf("_%s", collectorNumber)
	}

	if card.ImageURIs == nil {
		return nil, fmt.Errorf("no image URI's")
	}

	url := card.ImageURIs.Large

	urlParts := strings.Split(url, ".")
	extension := urlParts[len(urlParts)-1]
	if len(extension) > 5 {
		if !strings.Contains(extension, "?") {
			return nil, fmt.Errorf("strange file extension found, likely problem with logic of this program: %s", extension)
		}

		extension = strings.Split(extension, "?")[0]
	}

	imageData, err := download(url)
	if err != nil {
		return nil, fmt.Errorf("could not download: %v", err)
	}

	switch extension {
	case "png":
		if img, errDecode := png.Decode(bytes.NewReader(imageData)); errDecode == nil {
			images = append(images, addRoundedCornersWithThreshold(img, 0.5))
		}
	case "jpg", "jpeg":
		if img, errDecode := jpeg.Decode(bytes.NewReader(imageData)); errDecode == nil {
			// Create a new RGBA image of the same size as the decoded image
			bounds := img.Bounds()
			rgbaImg := image.NewRGBA(bounds)

			// Copy the pixels from the decoded image to the RGBA image
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					rgbaImg.Set(x, y, img.At(x, y))
				}
			}

			images = append(images, addRoundedCornersWithThreshold(rgbaImg, 0.5))
		}
	}

	return images, nil
}

func (s *Service) GetImagesFromDeckList(list string) (images []image.Image, err error) {
	cards := s.SearchWithDeckList(list)
	cards = s.scryfallGetFirstMatchCardsFromDeckList(list, cards)

	for _, card := range cards {
		cardImages, errGet := s.GetImagesFromCard(card)
		if errGet != nil {
			continue
		}

		if len(cardImages) < 1 {
			continue
		}

		images = append(images, cardImages[0])
	}

	return
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

func addRoundedCornersWithThreshold(img image.Image, normalizedThreshold float64) image.Image {
	bounds := img.Bounds()

	// Create a copy of the original image
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, img, image.Point{}, draw.Src)

	// Define the flood-fill starting points (corners)
	corners := []image.Point{
		{0, 0},                             // Top-left corner
		{bounds.Dx() - 1, 0},               // Top-right corner
		{0, bounds.Dy() - 1},               // Bottom-left corner
		{bounds.Dx() - 1, bounds.Dy() - 1}, // Bottom-right corner
	}

	// Flood-fill each corner
	for _, corner := range corners {
		floodFill(result, corner, normalizedThreshold)
	}

	return result
}

func floodFill(img *image.RGBA, start image.Point, normalizedThreshold float64) {
	bounds := img.Bounds()
	visited := make([][]bool, bounds.Dy())
	for y := range visited {
		visited[y] = make([]bool, bounds.Dx())
	}

	targetColor := img.At(start.X, start.Y).(color.RGBA)

	var queue []image.Point
	queue = append(queue, start)

	for len(queue) > 0 {
		p := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if visited[p.Y][p.X] {
			continue
		}

		visited[p.Y][p.X] = true

		currentColor := img.At(p.X, p.Y).(color.RGBA)
		distance := colorDistanceNormalized(targetColor, currentColor)

		if distance <= normalizedThreshold {
			img.Set(p.X, p.Y, color.Transparent)

			if p.X+1 < bounds.Dx() {
				queue = append(queue, image.Point{p.X + 1, p.Y})
			}
			if p.X-1 >= 0 {
				queue = append(queue, image.Point{p.X - 1, p.Y})
			}
			if p.Y+1 < bounds.Dy() {
				queue = append(queue, image.Point{p.X, p.Y + 1})
			}
			if p.Y-1 >= 0 {
				queue = append(queue, image.Point{p.X, p.Y - 1})
			}
		}
	}
}

func colorDistanceNormalized(a, b color.RGBA) float64 {
	// Calculate the squared Euclidean distance in RGB color space
	dR := float64(a.R) - float64(b.R)
	dG := float64(a.G) - float64(b.G)
	dB := float64(a.B) - float64(b.B)
	distanceSquared := dR*dR + dG*dG + dB*dB

	// Calculate the maximum possible distance (the diagonal of the RGB cube)
	maxDistance := math.Sqrt(3 * 255 * 255)

	// Normalize the distance by dividing by the maximum possible distance
	normalizedDistance := distanceSquared / (maxDistance * maxDistance)

	return normalizedDistance
}
