package mystem

import (
	"encoding/json"
	"strings"
)

type Grammeme int

// https://tech.yandex.ru/mystem/doc/grammemes-values-docpage/
const (
	// sex
	Male Grammeme = (iota + 1) << 1
	Female
	Neuter

	// todo add other grammems
	// part of speech
	// verb's time
	// case
	// number
	// verb inclination
	// form of adjectives
	// degree of comparison
	// verb's face
	// kind
	// voice
	// animation
	// transitivity
	// other notation
)

func (g Grammeme) Sex() Grammeme {
	if g&Male == Male {
		return Male
	} else if g&Female == Female {
		return Female
	} else if g&Neuter == Neuter {
		return Neuter
	}
	return 0
}

// --eng-gr required
func (g *Grammeme) UnmarshalJSON(data []byte) error {
	var (
		str         string
		newGrammeme Grammeme
	)

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	grammems := strings.Split(str, ",")
	for _, grammeme := range grammems {
		switch grammeme {
		// sex
		case "m":
			newGrammeme |= Male
		case "f":
			newGrammeme |= Female
		case "n":
			newGrammeme |= Neuter
		}
	}

	*g = newGrammeme
	return nil
}
