package DataTypes

import (
	"daphne/Helpers"
)

/**
 * Represents a page
 */
type Page struct {
	File       string
	OutFile    string
	Meta       map[string]string
	Content    []string
	IsBlogPost bool
}

func (self *Page) GetSlug() string {
	return Helpers.URLSafe(self.Meta["page.title"])
}

func (self *Page) GetPermalink(structure string) string {
	permalink := Helpers.Replace(structure, "%slug%", self.GetSlug())
	permalink = Helpers.Replace(permalink, "%year%", self.Meta["page.date_year"])
	permalink = Helpers.Replace(permalink, "%month%", self.Meta["page.date_month"])
	permalink = Helpers.Replace(permalink, "%day%", self.Meta["page.date_day"])

	return permalink
}

/**
 * An Inline command
 */
type InlineCommand struct {
	Control   string // Type of command (if, include, while, etc...)
	Condition string // Control condition
	IfTrue    string
	IfFalse   string
	StartLine int // Line the command starts at
	EndLine   int // Line the command ends at
}

/**
 * A multiline COmmand
 */
type MultilineCommand struct {
	Control   string // Type of command (if, include, while, etc...)
	Condition string // Control condition
	IfTrue    []string
	IfFalse   []string
	StartLine int // Line the command starts at
	EndLine   int // Line the command ends at
	State     int
}
