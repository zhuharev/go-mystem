package mystem_wrapper

import (
	"encoding/json"
	"strings"
)

type Grammeme int

// https://tech.yandex.ru/mystem/doc/grammemes-values-docpage/
const (
	// sex
	Male Grammeme = iota << 1
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

type Grammems struct {
	PartOfSpeech Grammeme
	Sex          Grammeme
}

// --eng-gr required
func (g *Grammems) UnmarshalJSON(data []byte) error {
	var str string

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	grammems := strings.Split(str, ",")
	for _, grammeme := range grammems {
		switch grammeme {
		// sex
		case "m":
			g.Sex = Male
		case "f":
			g.Sex = Female
		case "n":
			g.Sex = Neuter
		}
	}

	return nil
}
