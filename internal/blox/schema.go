package blox

import (
	_ "embed"
)

// ArticleCue is the cue for an Article
//go:embed article.cue
var ArticleCue string

// ArticleTemplate is the template for an article
//go:embed article.md
var ArticleTemplate string

// CategoryCue is the Cue for a category
//go:embed category.cue
var CategoryCue string

// CategoryTemplate is the template for a Category
//go:embed category.md
var CategoryTemplate string

// ProfileCue is the cue for a profile
//go:embed profile.cue
var ProfileCue string

// ProfileTemplate is the template for a profile
//go:embed profile.md
var ProfileTemplate string

// PageCue is the cue template for a page
//go:embed page.cue
var PageCue string

// PageTemplate is the page template
//go:embed page.md
var PageTemplate string
