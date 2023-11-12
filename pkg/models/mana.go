package models

type Mana int

const (
	ManaColorless Mana = iota
	ManaRed
	ManaBlue
	ManaBlack
	ManaWhite
	ManaGreen
	ManaPhyrexianColorless
	ManaPhyrexianRed
	ManaPhyrexianBlue
	ManaPhyrexianBlack
	ManaPhyrexianWhite
	ManaPhyrexianGreen
	NumManaTypes
)

func (m Mana) Name() string {
	return m.String()
}

func (m Mana) String() string {
	lookupTable := map[Mana]string{
		ManaColorless:          "Colorless",
		ManaRed:                "Red",
		ManaBlue:               "Blue",
		ManaBlack:              "Black",
		ManaWhite:              "White",
		ManaGreen:              "Green",
		ManaPhyrexianColorless: "Phyrexian Colorless",
		ManaPhyrexianRed:       "Phyrexian Red",
		ManaPhyrexianBlue:      "Phyrexian Blue",
		ManaPhyrexianBlack:     "Phyrexian Black",
		ManaPhyrexianWhite:     "Phyrexian White",
		ManaPhyrexianGreen:     "Phyrexian Green",
	}

	return lookupTable[m]
}
