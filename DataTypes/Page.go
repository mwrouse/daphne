package DataTypes

import (
    "daphne/Utils"
    "daphne/Constants"
)

type PageMeta map[string]string


/**
 * Represents a page
 */
type Page struct {
    File    string // The input file
    OutFile string // The output file

    Meta    PageMeta // Metadata for the page

    Content []string // Page content

    IsBlogPost bool
}


/**
 * Name.........: GetSlug
 * Return.......: string
 * Description..: Returns the slug of a blog post
 */
func (self Page) GetSlug() (string) {
    if self.IsBlogPost {
        return Utils.URLSafe(self.Meta[Constants.PAGE_TITLE])
    }

    return ""
}


/**
 * Name.........: GetPermalink
 * Parameters...: structure (string) - the permalink structure
 * Return.......: string
 * Description..: Gets the permalink structure for a post
 */
func (self Page) GetPermalink(structure string) (string) {
    permalink := Utils.Replace(structure, Constants.PERMALINK_SLUG, self.GetSlug())
    permalink = Utils.Replace(permalink, Constants.PERMALINK_YEAR, self.Meta[Constants.PAGE_YEAR])
    permalink = Utils.Replace(permalink, Constants.PERMALINK_MONTH, self.Meta[Constants.PAGE_MONTH])
    permalink = Utils.Replace(permalink, Constants.PERMALINK_DAY, self.Meta[Constants.PAGE_DAY])

    return permalink
}
