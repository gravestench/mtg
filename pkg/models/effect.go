package models

import (
	"fmt"
)

// EffectFlag represents a Magic: The Gathering card effect.
type EffectFlag int

const (
	// HasteEffect represents the Haste effect.
	HasteEffect EffectFlag = 1 << iota

	// FlyingEffect represents the Flying effect.
	FlyingEffect

	// LifelinkEffect represents the Lifelink effect.
	LifelinkEffect

	// TrampleEffect represents the Trample effect.
	TrampleEffect

	// DeathtouchEffect represents the Deathtouch effect.
	DeathtouchEffect

	// FirstStrikeEffect represents the First Strike effect.
	FirstStrikeEffect

	// DoubleStrikeEffect represents the Double Strike effect.
	DoubleStrikeEffect

	// VigilanceEffect represents the Vigilance effect.
	VigilanceEffect

	// HexproofEffect represents the Hexproof effect.
	HexproofEffect

	// IndestructibleEffect represents the Indestructible effect.
	IndestructibleEffect

	// MenaceEffect represents the Menace effect.
	MenaceEffect

	// ReachEffect represents the Reach effect.
	ReachEffect

	// ProwessEffect represents the Prowess effect.
	ProwessEffect
)

func (e EffectFlag) Names() (names []string) {
	lookupTable := map[EffectFlag]string{
		HasteEffect:          "Haste",
		FlyingEffect:         "Flying",
		LifelinkEffect:       "Lifelink",
		TrampleEffect:        "Trample",
		DeathtouchEffect:     "Deathtouch",
		FirstStrikeEffect:    "First",
		DoubleStrikeEffect:   "Double",
		VigilanceEffect:      "Vigilance",
		HexproofEffect:       "Hexproof",
		IndestructibleEffect: "Indestructible",
		MenaceEffect:         "Menace",
		ReachEffect:          "Reach",
		ProwessEffect:        "Prowess",
	}

	for _, t := range []EffectFlag{
		HasteEffect,
		FlyingEffect,
		LifelinkEffect,
		TrampleEffect,
		DeathtouchEffect,
		FirstStrikeEffect,
		DoubleStrikeEffect,
		VigilanceEffect,
		HexproofEffect,
		IndestructibleEffect,
		MenaceEffect,
		ReachEffect,
		ProwessEffect,
	} {
		if e&t > 0 {
			names = append(names, lookupTable[t])
		}
	}

	return
}

func (e EffectFlag) Descriptions() (descriptions []string) {
	lookupTable := map[EffectFlag]string{
		HasteEffect:          "Creatures with Haste can attack and tap on the turn they enter the battlefield.",
		FlyingEffect:         "Creatures with Flying can't be blocked except by creatures with Flying or Reach.",
		LifelinkEffect:       "When a creature with Lifelink deals damage, you gain that much life.",
		TrampleEffect:        "When a creature with Trample deals excess damage to a blocking creature, that damage is dealt to the defending player or planeswalker.",
		DeathtouchEffect:     "Any amount of damage dealt by a creature with Deathtouch is enough to destroy another creature.",
		FirstStrikeEffect:    "trike: Creatures with First Strike deal combat damage before creatures without it during combat.",
		DoubleStrikeEffect:   "trike: Creatures with Double Strike deal both first strike and regular combat damage during combat.",
		VigilanceEffect:      "Creatures with Vigilance don't tap when attacking and can block as normal.",
		HexproofEffect:       "A permanent with Hexproof can't be the target of spells or abilities your opponents control.",
		IndestructibleEffect: "Creatures and other permanents with Indestructible can't be destroyed by damage or effects that say 'destroy.'",
		MenaceEffect:         "Creatures with Menace can't be blocked except by two or more creatures.",
		ReachEffect:          "Creatures with Reach can block creatures with Flying as though they had Flying.",
		ProwessEffect:        "Whenever you cast a non-creature spell, creatures you control with Prowess get +1/+1 until the end of the turn.",
	}

	for _, t := range []EffectFlag{
		HasteEffect,
		FlyingEffect,
		LifelinkEffect,
		TrampleEffect,
		DeathtouchEffect,
		FirstStrikeEffect,
		DoubleStrikeEffect,
		VigilanceEffect,
		HexproofEffect,
		IndestructibleEffect,
		MenaceEffect,
		ReachEffect,
		ProwessEffect,
	} {
		if e&t > 0 {
			descriptions = append(descriptions, lookupTable[t])
		}
	}

	return descriptions
}

func (e EffectFlag) String() (result string) {
	names, descriptions := e.Names(), e.Descriptions()

	if len(names) < 1 {
		return ""
	}

	for index := range names {
		if result != "" {
			result = fmt.Sprintf("%s: %s", names[index], descriptions[index])
			continue
		}

		result = fmt.Sprintf("%s\r\n%s: %s", result, names[index], descriptions[index])
	}

	return
}
