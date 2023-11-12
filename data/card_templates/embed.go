package card_templates

import (
	_ "embed"
)

const (
	ArtworkOriginX = 0.5
	ArtworkOriginY = 0.334

	FrameLeft   = 88. / 743.
	FrameRight  = 656. / 743.
	FrameBottom = 567. / 1044.
	FrameTop    = 100. / 1044.
)

var (
	//go:embed default_artwork.png
	DefaultArtwork []byte

	//go:embed artifact.png
	Artifact []byte

	//go:embed back.png
	Back []byte

	//go:embed black.png
	Black []byte

	//go:embed blue.png
	Blue []byte

	//go:embed green.png
	Green []byte

	//go:embed land_artifact.png
	LandArtifact []byte

	//go:embed land_black.png
	LandBlack []byte

	//go:embed land_blue.png
	LandBlue []byte

	//go:embed land_green.png
	LandGreen []byte

	//go:embed land_red.png
	LandRed []byte

	//go:embed land_white.png
	LandWhite []byte

	//go:embed multicolor.png
	Multicolor []byte

	//go:embed original.png
	Original []byte

	//go:embed red.png
	Red []byte

	//go:embed white.png
	White []byte
)
