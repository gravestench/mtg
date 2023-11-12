package card

import (
	"image"

	"github.com/gravestench/mtg/pkg/models"
)

var _ cardInterface = &Card{}

type cardInterface interface {
	isCard
	hasManaCost
	hasAbilities
	canHaveCreatureFields
	hasUtilityMethods
	hasGraphics
}

type isCard interface {
	Tap() error
	Untap() error
	IsCardTapped() bool
	SuperType() models.CardSuperType
}

type hasManaCost interface {
	ManaCostString() string
	ConvertedManaCost() int
}

type hasAbilities interface {
	AddAbility(ability string)
	RemoveAbility(ability string)
}

type canHaveCreatureFields interface {
	Power() int
	Toughness() int
	SubTypes() []string
}

type hasUtilityMethods interface {
	CanTapOnFirstTurn() bool
}

type hasGraphics interface {
	Template() image.Image
	SetTemplate(image.Image)

	Artwork() image.Image
	SetArtwork(image.Image)

	CompositeCardImage() image.Image
}
