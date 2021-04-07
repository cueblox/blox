package markdown

import (
	"strings"
)

// ToYAML converts markdown to a YAML file
// storing the implicit 'body' of the markdown
// in 'body_raw'
func ToYAML(raw string) (string, error) {
	var content strings.Builder
	var err error
	lines := strings.Split(raw, "\n")
	var inBody bool
	for i, line := range lines {
		// remove first delimiter
		if i != 0 {
			if !inBody {
				// this is last delimiter
				// replace with 'body: |' and
				// indent the rest of the body by 2 spaces
				if line == "---" {
					content.WriteString("body_raw: |")
					content.WriteString("\n")
					inBody = true
				} else {
					content.WriteString(line)
					content.WriteString("\n")
				}
			} else {
				content.WriteString("  ")
				content.WriteString(line)
				content.WriteString("\n")
			}
		}

	}
	return content.String(), err
}

/*
--- first delimiter
key: value
key2: value
---
body text

===========

key: value
key2: value
body: |
  body text
  next line of body text
*/
