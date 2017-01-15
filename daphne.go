package main

import (
    "daphne/State"
    "daphne/Parser"
    "daphne/Helpers"
    "daphne/FileSystem"
    "daphne/DataTypes"
    "daphne/Errors"
    "fmt"
    "os"
    "io/ioutil"
    "regexp"
    "time"
)


var validFileExtensions, _ = regexp.Compile("^[^_\\.](.*)?$")

var ProgramState = State.NewCompilerState()

var fileWatch = make(map[string]time.Time)


/**
  * Name.........:
  * Parameters...:
  * Return.......:
  * Description..:
  */
func main() {
    wd, err := os.Getwd()
    if err != nil {
        (Errors.NewFatal(err.Error())).Handle()
    }
    ProgramState.Special = make(map[string][]DataTypes.Page)

    Helpers.Print("Green", "=======================================")
    Helpers.Print("Green", "       Daphne Static Website Builder   ")
    Helpers.Print("Green", "=======================================")
    fmt.Println("")
    fmt.Println("")

    argument := "build"
    if len(os.Args) > 0 {
        argument = Helpers.ToLower(os.Args[1])
    }


    switch (argument) {
    case "build":
        Build(wd)

    case "watch":
        Build(wd)
        Helpers.Print("White", "Monitoring...")
        for {
            if FileWatch(wd) {
                Build(wd)
                Helpers.Print("White", "\nMonitoring...")
            }
            time.Sleep(1000 * time.Millisecond)
        }

    case "new":
        dirs := []string{"_includes","_templates","_posts"}
        for _, fldr := range dirs {
            os.MkdirAll(fldr, 077)
        }
        
        if !FileSystem.FileExists("_config.daphne") {
            FileSystem.WriteFile("_config.daphne", []string{"site: {","}","blog: {", "}"})
        }

        Helpers.Print("Green", "Finished, default files have been created!")

    default:
        (Errors.NewFatal("Unknown argument: ", argument)).Handle()
    }
}






func Build(wd string) {
    Helpers.Print("White", "Building...")

    // Get information from the config file
    err := Parser.ParseConfigFile(wd, ProgramState)
    err.Handle()

    // Remove the output directory and its contents
    FileSystem.EmptyDir(ProgramState.Config["compiler.output"])

    // Loop through list of files
    Parser.PreparseFiles(ProgramState.Config["compiler.source"], ProgramState)

    // Expand Everything
    for _, pageType := range []string{"pages", "posts"} {
        for _, page := range ProgramState.GetSpecial("site." + pageType) {
            displayText := page.File
            if pageType == "posts" {
                displayText = page.Meta["page.title"]
            }

            Helpers.Print("Magenta", "\tBuilding: ", displayText)
            err := Parser.ExpandPage(&page, ProgramState)
            err.Handle()
        }
    }

    Helpers.Print("Green", "Finished")
}



func FileWatch(dir string) (bool) {
    list, err := ioutil.ReadDir(dir)
    if err != nil {
        (Errors.NewFatal(err.Error())).Handle()
    }

    for _, file := range list {
        name := ProgramState.Path(dir + "\\" + file.Name())

        if file.IsDir() {
            if file.Name() != ProgramState.Config["compiler.output"] {
                if FileWatch(dir + "\\" + file.Name()) {
                    return true
                }
            }
        } else {
            if !fileWatch[name].IsZero() {
                // Check if modifitcation time has changed
                if !fileWatch[name].Equal(file.ModTime()) {
                    Helpers.Print("Magenta", name + " was modified, will rebuild.")
                    fileWatch[name] = file.ModTime()
                    return true
                }
            } else {
                // Get modification time
                fileWatch[name] = file.ModTime()
            }
        }
    }

    return false
}
