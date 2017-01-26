package State

import (
    "daphne/DataTypes"
    "daphne/Utils"
    "daphne/Constants"
)


type CompilerState struct {
    Config      map[string]string
    Special     map[string]DataTypes.Page

    Ignore      []string // Files/Dirs to ignore

    CurrentPage DataTypes.Page // Current page being parsed
    MetaStack   DataTypes.Stack // Stack for the page meta (mainly used for for-loops)
}


/**
 * Name.........: NewCompilerState
 * Return.......: *CompilerState
 * Description..: Constructor for Compiler State
 */
func NewCompilerState() (*CompilerState) {
    state := new(CompilerState)

    state.Config = make(map[string]string)
    state.Ignore = []string{}
    state.MetaStack = make(DataTypes.DataStack, 0)

    return state
}


/**
 * Name.........: Exists
 * Parameters...: variable (string) - variable to look for
 * Return.......: bool
 * Description..: Check if a variable exists in conifg or current page meta
 */
func (self CompilerState) Exists(variable string) (bool) {
    currPage := self.MetaStack.Peek()

    return self.Config[variable] != "" || currPage[variable] != "" || self.CurrentPage.Meta[variable] != ""
}


/**
 * Name.........: Get
 * Parameters...: variable (string) - the variable to get
 * Return.......: string
 * Description..: Will get a variable
 */
func (self CompilerState) Get(variable string) (string) {
    currPage := self.MetaStack.Peek()

    // Meta stack has highest priority, then the current page, then the configuration of Daphne
    if currPage[variable] != "" {
        return currPage[variable]
    } else if self.CurrentPage.Meta[variable] != "" {
        return self.CurrentPage.Meta[variable]
    } else if self.Config[variable] != "" {
        return self.Config[variable]
    }

    return ""
}


/**
 * Name.........: Path
 * Parameters...: file (string) - the file to return the path to
 * Return.......: string
 * Description..: Returns a path to a file relative to the compiler source
 */
func (self CompilerState) Path(file string) (string) {
    path := self.Config[Constants.SOURCE_DIR] + "\\" + file
    path = Utils.Replace(path, "\\\\", "\\") // Replace double back slashes with one backslash
    path = Utils.Replace(path, ".\\", "") // Remove .\ from string

    return path
}


/**
 * Name.........: OutputPath
 * Parameters...: file (string) - the file to return a path to
 * Return.......: string
 * Description..: Returns the path of a file in relation to the output dir
 */
func (self CompilerState) OutputPath(file string) (string) {
    return self.Path(self.Config[Constants.OUTPUT_DIR] + "\\" + file)
}


/**
 * Name.........: Include
 * Parameters...: file (string) - file to return path to in includes dir
 * Return.......: string
 * Description..: Returns a path to a file in the includes dir
 */
func (self CompilerState) Include(file string) (string) {
    return self.Path(self.Config[Constants.INCLUDE_DIR] + "\\" + file)
}


/**
 * Name.........: Template
 * Parameters...: file (string) - file to get path for in the template directory
 * Return.......: string
 * Description..: Returns a path to a file in the template directory
 */
func (self CompilerState) Template(file string) (string) {
    return self.Path(self.Config[Constants.TEMPLATE_DIR] + "\\" + file)
}


/**
 * Name.........: IgnoreDir
 * Parameters...: dir (string) - the directory to check if it is ignore
 * Return.......: bool
 * Description..: Returns true if a directory should be ignored
 */
func (self CompilerState) IgnoreDir(dir string) (bool) {
    if len(dir) < 1 {
        return false // ?????
    }

    for _, fldr := range self.Ignore {
        // If dir is inside any of the folders that should be ignored, return true
        if Utils.IsInsideDir(dir, fldr) {
            return true
        }
    }

    return dir[:1] == "." && len(dir) > 1 // Ignore directories that start with a period
}


/**
 * Name.........: IgnoreDuringWatch
 * Parameters...: dir (string) - directory to check if ignoring
 * Return.......: bool
 * Description..: Determines if a directory should be ignored during watch
 */
func (self CompilerState) IgnoreDuringWatch(dir string) (bool) {
    if len(dir) < 1 {
        return false // ?????
    }

    // Only ignore stuff in a directory that starts with a period, or is inside the output dir
    return Utils.IsInsideDir(dir, self.Config[Constants.OUTPUT_DIR]) || (dir[:1] == "." && len(dir) > 1)
}


/**
 * Name.........: PageOutput
 * Parameters...: page (DataTypes.Page) - the page to get output path for
 * Return.......: string
 * Description..: Gets the output path for a page
 */
func (self CompilerState) PageOutput(page DataTypes.Page) (string) {
    path := Utils.Split(page.File, "\\") // Get directory parts

    // Remove .\ directory
    if path[0] == "." && len(path) > 1 {
        path = path[1:]
    }

    // Handle posts differently than other pages
    if path[0] == self.Config[Constants.POSTS_DIR] {
        // Is a blog post
        permalink := page.GetPermalink(self.Config[Constants.PERMALINK]) // Get the permalink for the page

        if self.Config[Constants.FOLDERICIZE] != "true" {
            // Do not move to own folder
            permalink = permalink + ".html"
        } else {
            // Put the blog post in its own folder
            permalink = permalink + "\\index.html"
        }

        return self.OutputPath(Utils.Replace(permalink, "/", "\\")) // Replace forward slashes with back slashes
    }

    // Not a blog post
    return self.OutputPath(page.File)
}


/**
 * Name.........: PageURL
 * Parameters...: page (DataTypes.Page) - the page to get the URL for
 * Return.......: string
 * Description..: Returns the URL to a page
 *
 * I used self.Path here instead of self.OutputPath because the output directory should
 * have the same directory structure as the source dir. If self.OutputPath was used then
 * all URLs would contain the output dir, which would not be good.
 */
func (self CompilerState) PageURL(page DataTypes.Page) (string) {
    // Handle the URL for blog posts differently
    if page.IsBlogPost {
        permalink := page.GetPermalink(self.Config[Constants.PERMALINK])

        if self.Config[Constants.FOLDERICIZE] != "true" {
            // Do not move to own folder
            permalink = permalink + ".html"
        } else {
            // Put the blog post in its own folder
            permalink = permalink + "\\index.html"
        }

        return self.Config[Constants.SITE_URL] + Utils.Replace(self.Path(permalink), "\\", "/")
    }

    // Not a blog post
    file := page.File
    filePath := Utils.Split(file, "\\")

    name := filePath[len(filePath) - 1] // Gets the file name
    filePath = filePath[:len(filePath) - 1] // Trim the file name

    extensionParts := Utils.Split(name, ".") // Break apart the file name
    extension := extensionParts[len(extensionParts) - 1]
    name = Utils.Join(extensionParts[:len(extensionParts) - 1], ".") // Get the name without the file extension

    path := utils.Join(filePath, "\\")

    url := ""

    if Utils.ToLower(name) == "index" {
        url = path + "\\" // Is a folder
    } else {
        url = path + "\\" + name + "." + extension // Is a file
    }

    return self.Config[Constants.SITE_URL] + Utils.Replace(self.Path(url), "\\", "/")
}
