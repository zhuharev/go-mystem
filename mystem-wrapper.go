package mystem_wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type MyStem struct {
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

func New(path string, args []string) *MyStem {
	m := &MyStem{}
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

func (m *MyStem) Transform(inputTexts []string) (transformedTexts []string, err error) {
	sentences, err := m.Sentences(inputTexts)
	if err != nil {
		return
	}

	for _, sentence := range sentences {
		sentenceText := ""
		for _, word := range sentence {
			if len(word.Analysis) != 0 {
				sentenceText = fmt.Sprintf("%s%s", sentenceText, word.Analysis[0].Lex)
			} else {
				sentenceText = fmt.Sprintf("%s%s", sentenceText, word.Text)
			}
		}

		transformedTexts = append(transformedTexts, sentenceText)
	}
	return
}

func (m *MyStem) Sentences(inputTexts []string) ([][]Word, error) {
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

	//parse every word into sentences
	var result [][]Word
	for i := range outByteTexts {
		var (
			words         []Word
			filteredWords []Word
		)

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

			filteredWords = append(filteredWords, words[wi])
		}

		result = append(result, filteredWords)
	}

	return result, nil
}
