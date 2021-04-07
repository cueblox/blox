package markdown

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	in := `---
key1: value
key2: value2
---
My Body
Body Line 2`
	expected := `key1: value
key2: value2
body_raw: |
  My Body
  Body Line 2
`
	output, err := ToYAML(in)
	if err != nil {
		t.Error(err)
	}
	if output != expected {
		fmt.Println(output)
		t.Error("output doesn't match expected")
	}
}
