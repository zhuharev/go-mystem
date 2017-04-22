package mystem_wrapper

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data := []byte(`"S,famn,f,anim=(acc,sg|gen,sg)"`)

	var gr Grammeme

	err := json.Unmarshal(data, &gr)
	if err != nil {
		t.Fatal(err)
	}
	if gr.Sex() != Female {
		t.Fatal("Sex not female")
	}
}
