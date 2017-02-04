package mystem_wrapper

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"encoding/json"
)

type myStem struct {
	path          string
	args          []string

	InputCharFilter *strings.Replacer
}

type word struct {
	Analysis []struct{
		Lex string `json:"lex"`
	} `json:"analysis"`
	Text string `json:"text"`
}

func New(path string, args []string) *myStem {
	m := &myStem{}
	m.path = path
	m.args = args

	// append c flag
	m.args = append(m.args, "-c", "--format", "json")

	// make replacer
	m.InputCharFilter = strings.NewReplacer(
		"«", "'",
		"»", "'",
		"—", "-",
		"_", "-",
		"\r\n", " ",
		"\r", " ",
		"\n", " ",
		"\t", " ",
	)

	return m
}

func (m *myStem) Transform(inputTexts []string) (transformedTexts []string, err error) {
	var inputBuffer, outBuffer bytes.Buffer

	for i := range inputTexts {
		// filter bad chars
		text := strings.TrimSpace(m.InputCharFilter.Replace(inputTexts[i]))
		if text == "" {
			text = " "
		}

		// one text = one line
		if i == len(inputTexts) - 1 {
			inputBuffer.Write([]byte(text))
		} else {
			inputBuffer.Write(append([]byte(text), '\n'))
		}
	}

	// run proc
	proc := exec.Command(m.path, m.args...)
	proc.Stdin = &inputBuffer
	proc.Stdout = &outBuffer
	err = proc.Start()
	if err != nil {
		return transformedTexts, fmt.Errorf("error running mystem: %s", err)
	}
	proc.Wait()
	if err != nil {
		return transformedTexts, fmt.Errorf("error waiting mystem: %s", err)
	}

	// parse output
	outByteTexts := bytes.Split(outBuffer.Bytes(), []byte("\n"))
	// remove always empty last string from mystem's output
	outByteTexts = outByteTexts[0:len(outByteTexts) - 1]
	if len(outByteTexts) != len(inputTexts) {
		return transformedTexts, fmt.Errorf("error: len(inputTexts)(%d) != len(outByteTexts)(%d) res: %s",
			len(inputTexts),
			len(outByteTexts),
			outBuffer.Bytes(),
		)
	}

	//parse every word
	for i := range outByteTexts {
		var words []word
		sentence := ""

		err := json.Unmarshal(outByteTexts[i], &words)
		if err != nil {
			return transformedTexts, fmt.Errorf("error while decoding json: %s res: %s",
				err,
				outByteTexts[i],
			)
		}

		for wi := range words {
			if words[wi].Text == "\n" {
				continue
			}

			if len(words[wi].Analysis) != 0 {
				sentence = fmt.Sprintf("%s%s", sentence, words[wi].Analysis[0].Lex)
			} else {
				sentence = fmt.Sprintf("%s%s", sentence, words[wi].Text)
			}
		}

		transformedTexts = append(transformedTexts, sentence)
	}

	return
}
