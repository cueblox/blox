package blox

import (
	"fmt"
)

func defaultTemplate(typ string) ([]byte, error) {

	switch typ {
	case "article":
		return []byte(ArticleTemplate), nil
	case "category":
		return []byte(CategoryTemplate), nil
	case "profile":
		return []byte(ProfileTemplate), nil
	case "page":
		return []byte(PageTemplate), nil
	default:
		return []byte{}, fmt.Errorf("generator doesn't support %s yet", typ)
	}
}
