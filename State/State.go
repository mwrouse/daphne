package State

import (
    "daphne/DataTypes"
    "daphne/Helpers"
)

/**
 * A struct to represent the current State
 */
type CompilerState struct {
    Config  map[string]string
    Special map[string][]DataTypes.Page

    CurrentPage DataTypes.Page
    Meta    DataTypes.MetaStack
}


/**
 * Compiler State Constructor
 */
func NewCompilerState() (*CompilerState) {
    state := new(CompilerState)

    state.Config = make(map[string]string)
    state.Meta = DataTypes.MetaStack{}

    return state;
}

/**
 * Returns true if a variable exists
 */
func (self CompilerState) Exists(variable string) (bool) {
    return self.Config[variable] != "" || self.CurrentPage.Meta[variable] != "" || (self.Meta.Peek())[variable] != ""
}


/**
 * Retrieves a variable from the compiler state
 */
func (self CompilerState) Get(variable string) (string) {
    if (self.Meta.Peek())[variable] != "" {
        return (self.Meta.Peek())[variable]
    } else if self.CurrentPage.Meta[variable] != "" {
        return self.CurrentPage.Meta[variable]
    } else if self.Config[variable] != "" {
        return self.Config[variable]
    }

    return ""
}


/**
 * Sets the value of a variable if possible
 */
func (self *CompilerState) Set(variable string, value string)  {
    if (*self).Config[variable] != "" {
        (*self).Config[variable] = value

    } else if ((*self).Meta.Peek())[variable] != "" {
        ((*self).Meta.Peek())[variable] = value
    }

    // Create New
    (*self).Config[variable] = value
}


/**
 * Gets a path in the compiler source
 */
func (self CompilerState) Path(file string) (string) {
    path := self.Config["compiler.source"] + "\\" + file

    return Helpers.Replace(Helpers.Replace(path, "\\\\", "\\"), ".\\", "")
}


/**
 * Gets the file path in the output
 */
func (self CompilerState) OutputPath(file string) (string) {
    return self.Path(self.Config["compiler.output"] + "\\" + file)
}


/**
 * Gets the path to a file in the _includes dir
 */
func (self CompilerState) Include(file string) (string) {
    return self.Path(self.Config["compiler.include_dir"] + "\\" + file)
}


/**
 * Gets the path a file in the templates dir
 */
func (self CompilerState) Template(file string) (string) {
    return self.Path(self.Config["compiler.template_dir"] + "\\" + file)
}


/**
 *
 */
func (self CompilerState) IgnoreDir(dir string) (bool) {
    if len(dir) < 1 {
        return false
    }

    if dir[:1] == "_"  {
        return dir == self.Config["compiler.include_dir"] || dir == self.Config["compiler.template_dir"] || dir == self.Config["compiler.output"]
    }
    return (dir[:1] == "." && len(dir) > 1)
}


func (self CompilerState) GetSpecial(name string) ([]DataTypes.Page) {
    return self.Special[name]
}


func (self CompilerState) GetPageOutpath(page DataTypes.Page) (string) {
    filePath := Helpers.Split(page.File, "\\")

    if filePath[0] == "." && len(filePath) > 1 {
        filePath = filePath[1:]
    }

    if filePath[0] == self.Config["compiler.posts_dir"] {
        page.IsBlogPost = true

        // Get Permalink structure
        permalink := page.GetPermalink(self.Config["blog.permalink"])

        if self.Config["blog.foldericize"] != "true" {
            permalink = permalink + ".html"
        } else {
            permalink = permalink + "\\index.html"
        }

        return self.OutputPath(Helpers.Replace(permalink, "/", "\\"))
    } else {
        return self.OutputPath(page.File)
    }
}


func (self CompilerState) GetPageURL(page DataTypes.Page) (string) {
    file := page.Meta["page.file"]
    filePath := Helpers.Split(file, "\\")

    name := filePath[len(filePath) - 1]
    filePath = filePath[:len(filePath) - 1]

    extensionParts := Helpers.Split(name, ".")
    extension := extensionParts[len(extensionParts) - 1]
    name = Helpers.Join(extensionParts[:len(extensionParts) -1 ], ".")

    path := Helpers.Join(filePath, "\\")

    url := ""
    if Helpers.ToLower(name) == "index" {
        url = path + "\\"
    } else {
        url = path + "\\" + name + "." + extension
    }

    return self.Config["site.url"] + Helpers.Replace(self.Path(url), "\\", "/")
}

func (self CompilerState) GetPostURL(page DataTypes.Page) (string) {
    permalink := page.GetPermalink(self.Config["blog.permalink"])

    if self.Config["blog.foldericize"] != "true" {
        permalink = permalink + ".html"
    } else {
        permalink = permalink + "\\"
    }

    return self.Config["site.url"] + Helpers.Replace(self.Path(permalink), "\\", "/")
}
