package Parser

import (
    "daphne/FileSystem"
    "daphne/State"
    "daphne/Helpers"
    "regexp"
)


var validFileExtensions, _ = regexp.Compile("^[^_\\.](.*)?$")


/**
  * Name.........: PreparseFiles
  * Parameters...: dir (string) - directory to get files in
  * Description..: Preparses (discovers) files
  */
func PreparseFiles(dir string, ProgramState *State.CompilerState) {
    files := FileSystem.CollapseDirectory(dir, "", true) // Get all the files in the directory

    // Loop through them
    for _, file := range files {
        if ProgramState.IgnoreDir(file.Directory) {
            continue
        }

        if validFileExtensions.MatchString(file.Info.Name()) {
            name := ProgramState.Path(file.Directory + "\\" + file.Info.Name())
            nameSplit := Helpers.Split(name, ".")
            ext := nameSplit[len(nameSplit) - 1]
            dest := ProgramState.OutputPath(name)

            if ext == "html" || ext == "htm" {
                // If in the posts directory then parse as a post
                if file.Directory == ProgramState.Config["compiler.posts_dir"] {
                    page, err := ParsePost(name, ProgramState)
                    err.Handle()
                    page.IsBlogPost = true

                    if !err.IsFatal() {
                        ProgramState.Special["site.posts"] = append(ProgramState.Special["site.posts"], page)
                    }
                } else if Helpers.Substring(file.Directory, 0, 0) == "_" {
                    dirPath := Helpers.Split(file.Directory, "\\")
                    varName := Helpers.Join(dirPath, ".")
                    varName = "site." + Helpers.Substring(varName, 1, len(varName) - 1)

                    page, err := ParsePage(name, ProgramState)
                    err.Handle()

                    // Copy into variable
                    ProgramState.Special[varName] = append(ProgramState.Special[varName], page)

                } else {
                    // Regular page
                    page, err := ParsePage(name, ProgramState)
                    err.Handle()

                    ProgramState.Special["site.pages"] = append(ProgramState.Special["site.pages"], page)
                }

            } else {
                // Just copy the file
                Helpers.Print("Cyan", "\tCopying: ", name, " => ", dest)
                err := FileSystem.CopyFile(name, dest)
                err.Handle()
            }
        }
    }
}
