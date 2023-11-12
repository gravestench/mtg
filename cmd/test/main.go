package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"

	"github.com/gravestench/mtg/data/card_templates"
	"github.com/gravestench/mtg/pkg/card"
	"github.com/gravestench/mtg/pkg/models"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	cost := card.ManaCost{
		models.ManaRed: 1,
	}

	c := card.New("Goblin", cost, 2, 2, models.HasteEffect, nil, models.Creature, "Goblin")

	fmt.Printf("Card Name: %v\r\n", c.Name)
	fmt.Printf("Mana Cost: %v\r\n", c.ManaCost)
	fmt.Printf("Effects: %v\r\n", c.State.Effects.Names())
	fmt.Printf("Converted Mana Cost: %v\r\n", c.ConvertedManaCost())

	template, err := loadPngFromBytes(card_templates.Blue)
	if err != nil {
		panic(fmt.Errorf("loading image: %v", err))
	}

	artwork, err := loadPngFromBytes(card_templates.DefaultArtwork)
	if err != nil {
		panic(fmt.Errorf("loading image: %v", err))
	}

	c.SetTemplate(template)
	c.SetArtwork(artwork)

	bounds := template.Bounds()
	width, height := int32(float64(int32(bounds.Max.X))*0.5), int32(float64(int32(bounds.Max.Y))*0.5)

	rl.InitWindow(width, height, "MTG")

	img := rl.NewImageFromImage(c.CompositeCardImage())
	imgRectangle := rl.NewRectangle(0, 0, float32(img.Width), float32(img.Height))
	scaledRectangle := rl.NewRectangle(0, 0, float32(img.Width/2), float32(img.Height/2))
	texture := rl.LoadTextureFromImage(img)

	var r float32

	for !rl.WindowShouldClose() {
		r = float32(math.Sin(float64(time.Now().UnixNano())/500000000.0) * 1)
		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{0, 0, 0, 255})
		rl.DrawTexturePro(texture, imgRectangle, scaledRectangle, rl.Vector2{X: 0.5, Y: 0.5}, r, color.RGBA{255, 255, 255, 255})
		rl.EndDrawing()
	}
}

func loadPngFromBytes(data []byte) (image.Image, error) {
	buf := bytes.NewReader(data)

	img, err := png.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("decoding png: %v", err)
	}

	return img, nil
}
