package mystem_wrapper

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data := []byte(`"S,famn,m,anim=(acc,sg|gen,sg)"`)

	gr := Grammeme{}

	err := json.Unmarshal(data, &gr)
	if err != nil {
		t.Fatal(err)
	}
	if gr.Sex != Male {
		t.Fatal("Sex not male")
	}
}
