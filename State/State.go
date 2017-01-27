package State
/*
import (
    "daphne/DataTypes"
    "daphne/Helpers"
)

type SpecialFunction func(DataTypes.Page, *CompilerState2)


/**
 * A struct to represent the current State
 *
type CompilerState2 struct {
    Config  map[string]string
    Special map[string][]DataTypes.Page
    Ignore []string

    CurrentPage DataTypes.Page
    Meta    DataTypes.Stack

    PerformAfterFileWrite []SpecialFunction
}

/**
 * Compiler State Constructor
 *
func NewCompilerState2() (*CompilerState2) {
    state := new(CompilerState2)

    state.Config = make(map[string]string)
    state.Meta = make(DataTypes.Stack, 0)
    state.Ignore = []string{}

    state.PerformAfterFileWrite = []SpecialFunction{}

    return state;
}

/**
 * Returns true if a variable exists
 *
func (self CompilerState2) Exists(variable string) (bool) {
    return self.Config[variable] != "" || self.CurrentPage.Meta[variable] != "" || (self.Meta.Peek())[variable] != ""
}


/**
 * Retrieves a variable from the compiler state
 *
func (self CompilerState2) Get(variable string) (string) {
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
 *
func (self *CompilerState2) Set(variable string, value string)  {
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
 *
func (self CompilerState2) Path(file string) (string) {
    path := self.Config["compiler.source"] + "\\" + file

    return Helpers.Replace(Helpers.Replace(path, "\\\\", "\\"), ".\\", "")
}


/**
 * Gets the file path in the output
 *
func (self CompilerState2) OutputPath(file string) (string) {
    return self.Path(self.Config["compiler.output"] + "\\" + file)
}


/**
 * Gets the path to a file in the _includes dir
 *
func (self CompilerState2) Include(file string) (string) {
    return self.Path(self.Config["compiler.include_dir"] + "\\" + file)
}


/**
 * Gets the path a file in the templates dir
 *
func (self CompilerState2) Template(file string) (string) {
    return self.Path(self.Config["compiler.template_dir"] + "\\" + file)
}


/**
 *
 *
func (self CompilerState2) IgnoreDir(dir string) (bool) {
    if len(dir) < 1 {
        return false
    }

    //dirPath := Helpers.Split(dir, "\\")

    for _, fldr := range self.Ignore {
        if Helpers.IsInsideDir(dir, fldr) {
            return true
        }
    }
    /*
    if dir[:1] == "_"  {
        dir2 := dirPath[0]

        for _, fldr := range self.Ignore {
            if Helpers.HasDir(dir, fldr) {
                return true;
            }
        }*

    return (dir[:1] == "." && len(dir) > 1)
}


func (self CompilerState2) IgnoreDirDuringWatch(dir string) (bool) {
    if len(dir) < 1 {
        return false
    }

    dirPath := Helpers.Split(dir, "\\")

    if dir[:1] == "_"  {
        dir = dirPath[0]

        return dir == self.Config["compiler.output"]
    }
    return (dir[:1] == "." && len(dir) > 1)
}



func (self CompilerState2) GetSpecial(name string) ([]DataTypes.Page) {
    return self.Special[name]
}


func (self CompilerState2) GetPageOutpath(page DataTypes.Page) (string) {
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


func (self CompilerState2) GetPageURL(page DataTypes.Page) (string) {
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

func (self CompilerState2) GetPostURL(page DataTypes.Page) (string) {
    permalink := page.GetPermalink(self.Config["blog.permalink"])

    if self.Config["blog.foldericize"] != "true" {
        permalink = permalink + ".html"
    } else {
        permalink = permalink + "\\"
    }

    return self.Config["site.url"] + Helpers.Replace(self.Path(permalink), "\\", "/")
}
*/
