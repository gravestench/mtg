package card

import (
	"github.com/gravestench/mtg/pkg/models"
)

type CardState struct {
	IsTapped bool
	Effects  models.EffectFlag
}
