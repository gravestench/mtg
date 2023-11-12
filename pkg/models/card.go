package models

type CardSuperType int

const (
	Artifact CardSuperType = iota
	Creature
	Spell
	Instant
	Sorcery
	Enchantment
	Land
)
