package card

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"

	"github.com/gravestench/mtg/data/card_templates"
)

func (c *Card) Template() image.Image {
	if c.Graphics.Template != nil {
		return c.Graphics.Template
	}

	img, err := png.Decode(bytes.NewReader(card_templates.Artifact))
	if err != nil {
		return nil
	}

	return img
}

func (c *Card) SetTemplate(img image.Image) {
	if img == nil {
		img = c.Template()
	}

	c.Graphics.Template = img
}

func (c *Card) Artwork() image.Image {
	if c.Graphics.Artwork != nil {
		return c.Graphics.Artwork
	}

	data := bytes.NewReader(card_templates.DefaultArtwork)
	img, err := png.Decode(data)
	if err != nil {
		return nil
	}

	return img
}

func (c *Card) SetArtwork(img image.Image) {
	if img == nil {
		img = c.Artwork()
	}

	c.Graphics.Artwork = img
}

func (c *Card) CompositeCardImage() image.Image {
	artwork, template := c.Artwork(), c.Template()
	// Create a new RGBA image with the same size as your input images
	artworkBounds := artwork.Bounds()
	templateBounds := template.Bounds()
	result := image.NewRGBA(templateBounds)

	frameCenter := image.Point{
		X: int(float64(templateBounds.Dx()) * card_templates.ArtworkOriginX),
		Y: int(float64(templateBounds.Dy()) * card_templates.ArtworkOriginY),
	}

	origin := image.Point{
		X: frameCenter.X - (artworkBounds.Dx() / 2),
		Y: frameCenter.Y - (artworkBounds.Dy() / 2) + 40,
	}

	frameRectangle := image.Rectangle{
		Min: image.Point{
			X: int(float64(template.Bounds().Dx()) * card_templates.FrameLeft),
			Y: int(float64(template.Bounds().Dy()) * card_templates.FrameTop),
		},
		Max: image.Point{
			X: int(float64(template.Bounds().Dx()) * card_templates.FrameRight),
			Y: int(float64(template.Bounds().Dy()) * card_templates.FrameBottom),
		},
	}

	// Composite the first image onto the result image
	draw.Draw(result, frameRectangle, artwork, origin, draw.Over)
	origin.Y -= 60
	draw.Draw(result, frameRectangle, artwork, origin, draw.Over)

	// Composite the second image onto the result image, taking into account transparency
	draw.Draw(result, templateBounds, template, image.Point{}, draw.Over)

	return result
}
