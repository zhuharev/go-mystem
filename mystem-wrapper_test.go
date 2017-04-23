package mystem_wrapper

import (
	"os"
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {
	inputs := []string{
		`Данный пакет позволяет запускать Mystem из golang'a и получать выходной текст,`,
		`где каждое слово приведено в свою исходную форму.`,
		`Это нужно для подготовки датасетов для различных семантических анализов.`,
	}
	expected := []string{
		`данный пакет позволять запускать Mystem из golang'a и получать выходной текст,`,
		`где каждый слово приводить в свой исходный форма.`,
		`это нужно для подготовка датасет для различный семантический анализ.`,
	}

	mystem := New(os.Getenv("MYSTEM_PATH"), []string{"-d"})
	outputs, err := mystem.Transform(inputs)
	if err != nil {
		t.Fatal(err)
	}
	for i, got := range outputs {
		if strings.TrimSpace(got) != expected[i] {
			t.Fatalf("Got: %s != %s (expected)", got, expected[i])
		}
	}
}
