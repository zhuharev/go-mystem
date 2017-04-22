package mystem

import (
	"encoding/json"
	"strings"
)

type Grammeme struct {
	PartOfSpeech string
	Sex          Sex
}

// --eng-gr required
func (g *Grammeme) UnmarshalJSON(data []byte) error {
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

type Sex string

const (
	Male   Sex = "мужской"
	Female     = "женский"
	Neuter     = "средний"
)
