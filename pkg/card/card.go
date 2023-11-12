package card

import (
	"errors"
	"fmt"
	"image"
	"strings"

	"github.com/gravestench/mtg/pkg/models"
)

type ManaCost map[models.Mana]int

// New creates a new Card instance
func New(name string, cost map[models.Mana]int, power, toughness int, effects models.EffectFlag, abilities []string, superType models.CardSuperType, subTypes ...string) *Card {
	c := Builder().
		Name(name).
		ManaCost(cost).
		Power(power).
		Effects(effects).
		Toughness(toughness).
		SuperType(superType).
		SubTypes(subTypes).
		Build()

	for _, ability := range abilities {
		c.AddAbility(ability)
	}

	return c
}

// Card represents a game card
type Card struct {
	Name string
	ManaCost
	IsPermanent bool

	State CardState

	power     int
	toughness int

	Counters Counters

	Abilities map[string]any

	superType models.CardSuperType
	subTypes  []string

	Graphics struct {
		Template image.Image
		Artwork  image.Image
	}
}

func (c *Card) ManaCostString() string {
	var list []string

	for t := models.Mana(0); t < models.NumManaTypes; t++ {
		if cost, exists := c.ManaCost[t]; exists {
			list = append(list, fmt.Sprintf("%s: %d", t.String(), cost))
		}
	}

	return strings.Join(list, ", ")
}

// Tap the card
func (c *Card) Tap() error {
	if c.State.IsTapped {
		return errors.New("already tapped")
	}

	c.State.IsTapped = true

	return nil
}

// Untap the card
func (c *Card) Untap() error {
	if !c.State.IsTapped {
		return errors.New("already untapped")
	}

	c.State.IsTapped = false

	return nil
}

// IsCardTapped checks if the card is tapped
func (c *Card) IsCardTapped() bool {
	return c.State.IsTapped
}

// ConvertedManaCost calculates the converted mana cost of the card
func (c *Card) ConvertedManaCost() int {
	cmc := 0

	for _, cost := range c.ManaCost {
		cmc += cost
	}

	return cmc
}

// Power returns the card's power (for creatures)
func (c *Card) Power() int {
	if c.superType != models.Creature {
		return 0
	}

	total := c.power

	for _, counter := range c.Counters.Power {
		total += counter
	}

	return total
}

// Toughness returns the card's toughness (for creatures)
func (c *Card) Toughness() int {
	if c.superType != models.Creature {
		return 0
	}

	total := c.toughness

	for _, counter := range c.Counters.Toughness {
		total += counter
	}

	return total
}

// AddAbility adds an ability to the card
func (c *Card) AddAbility(ability string) {
	c.Abilities[ability] = nil
}

// RemoveAbility removes an ability from the card
func (c *Card) RemoveAbility(ability string) {
	delete(c.Abilities, ability)
}

func (c *Card) SuperType() models.CardSuperType {
	return c.superType
}

func (c *Card) SubTypes() []string {
	return c.subTypes
}

func (c *Card) CanTapOnFirstTurn() bool {
	if (c.State.Effects & models.HasteEffect) > 0 {
		return true
	}

	if c.superType == models.Land {
		return true
	}

	return false
}
