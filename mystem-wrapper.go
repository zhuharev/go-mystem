package mystem_wrapper

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	TEXT_SEPARATOR  = "~!=~#~=!~"
	REGEX_EXTRACTOR = `(.*?)\{(.+?)\|?\?{0,2}\}([\s\.\?\+\-\!\,\"\'\#\$\%\&\@\*\(\)\[\]\\\/\=\;\~\_]+)`
)

type myStem struct {
	path          string
	args          []string
	wordExtractor *regexp.Regexp

	InputCharFilter *strings.Replacer
}

func New(path string, args []string) *myStem {
	m := &myStem{}
	m.path = path
	m.args = args

	// append c flag
	m.args = append(m.args, "-c")

	// make replacer
	m.InputCharFilter = strings.NewReplacer(
		"{", "",
		"}", "",
		"«", "'",
		"»", "'",
		"—", "-",
		"_", "-",
		"\r\n", " ",
		"\n", " ",
		"\t", " ",
		TEXT_SEPARATOR, "---",
	)
	// make regex to take the first initial word form
	m.wordExtractor = regexp.MustCompile(REGEX_EXTRACTOR)

	return m
}

func (m *myStem) Transform(inputTexts []string) (transformedTexts []string, err error) {
	var inputBuffer, outBuffer bytes.Buffer

	// form input buffer
	byteSeparator := []byte(TEXT_SEPARATOR)
	for i := range inputTexts {
		// filter bad chars
		text := m.InputCharFilter.Replace(inputTexts[i])

		// include byte separator
		if i != 0 {
			inputBuffer.Write(byteSeparator)
		}

		// one text = one line
		inputBuffer.Write(append([]byte(text), '\n'))
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
	outBytes := bytes.Split(outBuffer.Bytes(), byteSeparator)
	for i := range outBytes {
		transformedTexts = append(
			transformedTexts,
			string(m.wordExtractor.ReplaceAll(bytes.Trim(outBytes[i], "\n"), []byte("$2$3"))),
		)
	}

	// ensure that count of input texts exacts count of transformed texts
	if len(inputTexts) != len(transformedTexts) {
		return transformedTexts, fmt.Errorf("error: len(inputTexts)(%d) != len(transformedTexts)(%d) res: %s",
			len(inputTexts),
			len(transformedTexts),
			outBytes,
		)
	}

	return
}
