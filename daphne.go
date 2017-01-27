package main

import (
    "daphne/State"
    "daphne/Parser"
    "daphne/Helpers"
    "daphne/FileSystem"
    "daphne/DataTypes"
    "daphne/Errors"
    "bufio"
    "fmt"
    "os"
    "regexp"
    "time"
    "net/http"
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
    if len(os.Args) > 1 {
        argument = Helpers.ToLower(Helpers.Join(os.Args[1:len(os.Args)], " "))
    }
    argument = Helpers.Trim(argument)

    // Pre-build on everything except help
    if argument != "help" {
        PreBuild(wd)
    }

    // Perform actions based on the argument
    switch (argument) {
    case "build":
        Build(wd)

    case "watch":
        Watch(wd)

    case "new":
        NewProject()

    case "serve":
        Serve(wd)

    case "new post":
        NewPost()

    case "help":
        Helpers.Print("white", "Arguments:")
        Helpers.Print("white", "\tbuild     - Build your website")
        Helpers.Print("white", "\tnew       - Create basic folder structure for new projects")
        Helpers.Print("white", "\tnew post  - Create a new post")
        Helpers.Print("white", "\twatch     - Watch for file changes and build when they are changed")
        Helpers.Print("white", "\tserve     - Host website on local web server, and watch for changes")
        fmt.Println("")

    default:
        (Errors.NewFatal("Unknown argument: ", argument)).Handle()
    }
}


/**
 * Name.........: PreBuild
 * Parameters...: wd (string) - the working directory
 * Description..: Discovers files and reads the Daphne configuration
 */
func PreBuild(wd string) {
    Helpers.Print("White", "Pre-Build...")

    // Get config file info
    err := Parser.ParseConfigFile(wd, ProgramState)
    err.Handle()

    // Clear the output directory
    FileSystem.EmptyDir(ProgramState.Config["compiler.output"])
}


/**
 * Name.........: Build
 * Parameters...: wd (string)
 * Description..: Builds all of the files (runs after PreBuild)
 */
func Build(wd string) {
    // Preparse all the files so they can reference one another
    Parser.PreparseFiles(ProgramState.Config["compiler.source"], ProgramState)

    Helpers.Print("White", "Building...")

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


/**
 * Name.........: Watch
 * Parameters...: wd (string)
 * Description..: Watches for file changes in an infinite loop
 */
func Watch(wd string) {
    Build(wd)
    Helpers.Print("White", "Monitoring...")
    for {
        // Watch for file changes, build if one has changed
        if FileWatch(wd) {
            ProgramState.Special = make(map[string][]DataTypes.Page) // Clear to avoid duplicates
            Build(wd)
        }
        time.Sleep(1000 * time.Millisecond)
    }
}


/**
 * Name.........: NewProject
 * Description..: Creates basic file structure
 */
func NewProject() {
    dirs := []string{"_includes","_templates","_posts"}
    for _, fldr := range dirs {
        os.MkdirAll(fldr, 077)
    }

    if !FileSystem.FileExists("_config.daphne") {
        FileSystem.WriteFile("_config.daphne", []string{"site: {","}","blog: {", "}"})
    }

    Helpers.Print("Green", "Finished, default files have been created!")
}


/**
 * Name.........: NewPost
 * Description..: Creates a new post
 */
func NewPost() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Post Title: ")

    title, _ := reader.ReadString('\n')
    title = Helpers.Replace(title, "\n", "")

    t := time.Now()

    path := ProgramState.Config["compiler.posts_dir"] + "\\" + t.Format("2006-01-02") + "-" + Helpers.URLSafe(title) + ".html"
    images := ProgramState.Config["compiler.posts_image_dir"] + "\\" + Helpers.URLSafe(title)

    // Create the post file and the directory
    FileSystem.WriteFile(path, []string{"---", "title: " + title, "template: post", "---"})
    FileSystem.CreateDir(images)


}


/**
 * Name.........: Serve
 * Parameters...: wd (string)
 * Description..: Builds files and starts a web server
 */
func Serve(wd string) {
    // Modify the URL of the website
    ProgramState.Config["site.url"] = "http://localhost:8081/"

    // Start the web server
    http.Handle("/", http.FileServer(http.Dir("./" + ProgramState.Config["compiler.output"])))

    Helpers.Print("Yellow", "\n\nWeb Server Started: ", ProgramState.Config["site.url"])
    go http.ListenAndServe(":8081", nil)
    Watch(wd)
}


/**
 * Name.........: FileWatch
 * Parameters...: dir (string)
 * Return.......: bool - true if a file has changed
 * Description..: Monitors files for changes
 */
func FileWatch(dir string) (bool) {
    files := FileSystem.CollapseDirectory(dir, "", true) // Get all the files in the directory

    for _, file := range files {
        if ProgramState.IgnoreDuringWatch(file.Directory) {
            continue
        }

        if validFileExtensions.MatchString(file.Info.Name()) {
            name := ProgramState.Path(file.Directory + "\\" + file.Info.Name())

            if !fileWatch[name].IsZero() {
                // Check if modification time is different
                if !fileWatch[name].Equal(file.Info.ModTime()) {
                    fmt.Println(file.Directory)

                    Helpers.Print("Magenta", name + " was modified, will rebuild.")
                    fileWatch[name] = file.Info.ModTime()
                    return true
                }
            } else {
                fileWatch[name] = file.Info.ModTime()
            }
        }
    }

    return false
}
