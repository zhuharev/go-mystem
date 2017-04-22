package mystem_wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type myStem struct {
	path string
	args []string

	InputCharFilter *strings.Replacer
}

type Word struct {
	Analysis []struct {
		Lex      string   `json:"lex"`
		Grammeme Grammeme `json:"gr"`
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
	words, err := m.Words(inputTexts)
	if err != nil {
		return
	}
	for _, word := range words {
		sentence := ""
		if len(word.Analysis) != 0 {
			sentence = fmt.Sprintf("%s%s", sentence, word.Analysis[0].Lex)
		} else {
			sentence = fmt.Sprintf("%s%s", sentence, word.Text)
		}
		transformedTexts = append(transformedTexts, sentence)
	}
	return
}

func (m *myStem) Words(inputTexts []string) ([]Word, error) {
	var inputBuffer, outBuffer bytes.Buffer

	for i := range inputTexts {
		// filter bad chars
		text := strings.TrimSpace(m.InputCharFilter.Replace(inputTexts[i]))
		if text == "" {
			text = " "
		}

		// one text = one line
		if i == len(inputTexts)-1 {
			inputBuffer.Write([]byte(text))
		} else {
			inputBuffer.Write(append([]byte(text), '\n'))
		}
	}

	// run proc
	proc := exec.Command(m.path, m.args...)
	proc.Stdin = &inputBuffer
	proc.Stdout = &outBuffer
	err := proc.Start()
	if err != nil {
		return nil, fmt.Errorf("error running mystem: %s", err)
	}
	proc.Wait()
	if err != nil {
		return nil, fmt.Errorf("error waiting mystem: %s", err)
	}

	// parse output
	outByteTexts := bytes.Split(outBuffer.Bytes(), []byte("\n"))
	// remove always empty last string from mystem's output
	outByteTexts = outByteTexts[0 : len(outByteTexts)-1]
	if len(outByteTexts) != len(inputTexts) {
		return nil, fmt.Errorf("error: len(inputTexts)(%d) != len(outByteTexts)(%d) res: %s",
			len(inputTexts),
			len(outByteTexts),
			outBuffer.Bytes(),
		)
	}

	var result []Word

	//parse every word
	for i := range outByteTexts {
		var words []Word

		err := json.Unmarshal(outByteTexts[i], &words)
		if err != nil {
			return nil, fmt.Errorf("error while decoding json: %s res: %s",
				err,
				outByteTexts[i],
			)
		}

		for wi := range words {
			if words[wi].Text == "\n" {
				continue
			}

			result = append(result, words[wi])
		}
	}

	return result, nil
}
