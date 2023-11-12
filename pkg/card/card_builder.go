package card

import (
	"github.com/gravestench/mtg/pkg/models"
)

func Builder() *CardBuilder {
	return &CardBuilder{
		name:        "Name",
		manaCost:    ManaCost{models.ManaColorless: 1},
		isPermanent: true,
		power:       1,
		toughness:   1,
		abilities:   make(map[string]any),
		superType:   models.Creature,
		subTypes:    make([]string, 0),
	}
}

type CardBuilder struct {
	name        string
	manaCost    ManaCost
	isPermanent bool
	effects     models.EffectFlag
	power       int
	toughness   int
	abilities   map[string]any
	superType   models.CardSuperType
	subTypes    []string
}

func (c *CardBuilder) Build() *Card {
	return &Card{
		Name:        c.name,
		ManaCost:    c.manaCost,
		IsPermanent: c.isPermanent,
		State: CardState{
			Effects: c.effects,
		},
		power:     c.power,
		toughness: c.toughness,
		Counters: Counters{
			Power:     make([]int, 0),
			Toughness: make([]int, 0),
		},
		Abilities: c.abilities,
		superType: c.superType,
		subTypes:  c.subTypes,
	}
}

func (c *CardBuilder) Name(s string) *CardBuilder {
	c.name = s

	return c
}

func (c *CardBuilder) ManaCost(m map[models.Mana]int) *CardBuilder {
	if m != nil {
		c.manaCost = m
	}

	return c
}

func (c *CardBuilder) IsPermanent(b bool) *CardBuilder {
	c.isPermanent = b
	return c
}

func (c *CardBuilder) Power(i int) *CardBuilder {
	c.power = i
	return c
}

func (c *CardBuilder) Toughness(i int) *CardBuilder {
	c.toughness = i
	return c
}

func (c *CardBuilder) Effects(e models.EffectFlag) *CardBuilder {
	c.effects = e

	return c
}

func (c *CardBuilder) Abilities(m map[string]any) *CardBuilder {
	if m != nil {
		c.abilities = m
	}

	return c
}

func (c *CardBuilder) SuperType(superType models.CardSuperType) *CardBuilder {
	c.superType = superType

	return c
}

func (c *CardBuilder) SubTypes(s []string) *CardBuilder {
	if s != nil {
		c.subTypes = s
	}

	return c
}
