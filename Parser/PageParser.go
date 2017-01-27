package Parser


import (
    "daphne/FileSystem"
    . "daphne/DataTypes"
    "daphne/State"
    "daphne/Helpers"
    "daphne/Grammar"
    "daphne/Grammar/Semantics"
    "daphne/Errors"
    "time"
)



/**
 * Name.........: ParsePage
 * Parameters...: file (string) - name of file to parse
 *                ProgramState (*State.CompilerState) - The State
 * Return.......: *Page - the parsed page structure
 *                error - any errors
 * Description..: Parses a page into the meta section and content section
 */
func ParsePage(file string, ProgramState *State.CompilerState) (Page, Errors.Error) {
    contents, err := FileSystem.ReadFile(file)
    if err.HasError() {
        return Page{}, err
    }

    // Create struct to hold the page information
    page := Page{Meta:make(map[string]string), Content:[]string{}, File: file}

    state := 0

    /**
    * State machine for splitting the config section and the content section
    */
    for i, origLine := range contents {
        line := Helpers.Trim(origLine)

        switch (state) {
        // 0 - looking for config opener
        case 0:
            if line != ProgramState.Config["compiler.tags.meta"] {
                return page, Errors.NewFatal("First line of ", file, " can ONLY Be the opening meta tags")
            }

            state = 1

        // 1 - Waiting for end meta config
        case 1:
            if line == ProgramState.Config["compiler.tags.meta"] {
                meta := Helpers.Copy(contents[:i])
                meta[0] = "page: {"
                meta = append(meta, "}")

                contents[i] = "" // Clear the ending config line

                // Parse the config data read so far
                err := Errors.None()
                page.Meta, err = ParseConfig(meta)
                if err.HasError() {
                    return page, err
                }

                page.Meta["page.file"] = file
                page.Meta["page.url"] = ProgramState.PageURL(page)

                ApplyDefaultMetaConfig(page.Meta, ProgramState)

                state = 2 // Advance to the next state
            }

        // Wait for something that is not a blank line (removes lots of empty space between config and html)
        case 2:
            if line != "" {
                page.Content = Helpers.Copy(contents[i:])
                page.Meta["content"] = Helpers.Join(page.Content, "\n")

                page.OutFile = ProgramState.PageOutput(page)

                GetExcerpt(&page, ProgramState)
                return page, Errors.None()
            }
        }
    }

    // If the state is 2, then the file is blank, let this be okay
    if state == 2 {
        page.OutFile = ProgramState.PageOutput(page)
        return page, Errors.NewWarning(file, " is empty\n")
    }

    // If execution reaches here there was no end configuration, or it was invalid and not on a separate line
    return page, Errors.NewFatal("No end to config section was found for: ", file)
}


/**
 * Name.........: ParsePost
 * Parameters...:
 * Return.......:
 * Description..:
 */
func ParsePost(file string, ProgramState *State.CompilerState) (Page, Errors.Error) {
    page, err := ParsePage(file, ProgramState) // Parse it as a page
    if err.HasError() {
        return page, err
    }

    // Get the file name from the file path
    nameParts := Helpers.Split(file, "\\")
    name := nameParts[len(nameParts) - 1]

    page.IsBlogPost = true

    // Get the post date
    date := Helpers.Split(name, "-") // Split at dashes
    if len(date) < 4 {
        return page, Errors.NewWarning("The file ", name, " has an invalid format, it should be YYYY-MM-DD-identifier.html, it will not be added to the list of posts")
    }

    t, err1 := time.Parse("2006-01-02", Helpers.Join(date[:3], "-")) // Parse the date
    if err1 != nil {
        return page, Errors.NewFatal(err1.Error())
    }

    page.Meta["page.date"] = t.Format("January 2, 2006")
    page.Meta["page.date_year"] = t.Format("2006")
    page.Meta["page.date_month"] = t.Format("01")
    page.Meta["page.date_day"] = t.Format("02")
    page.Meta["page.date_short"] = t.Format("2006-01-02")

    // Get the post URL Information
    page.Meta["page.slug"] = page.GetSlug()
    page.Meta["page.url"] = ProgramState.PageURL(page)

    page.OutFile = ProgramState.PageOutput(page)

    return page, Errors.None()
}


/**
 * Name.........: ApplyDefaultMetaConfig
 * Parameters...: meta (map[string]string) - the config to apply defualts to
 *                ProgarmState (*State.CompilerState) - Compiler state
 * Description..: Adds default options to a configuration if they do not exist
 */
func ApplyDefaultMetaConfig(meta map[string]string, ProgramState *State.CompilerState) {
    defaults := []string{"author", "description", "title", "template"}

    for _, key := range defaults {
        if meta["page." + key] == "" {
            meta["page." + key] = ProgramState.Config["site." + key]
        }
    }
}


/**
 * Name.........:
 * Parameters...:
 * Return.......:
 * Description..:
 */
func GetExcerpt(page *Page, ProgramState *State.CompilerState) {
    if page.Meta["page.excerpt"] != "" {
        return
    }

    excerpt := []string{}

    for _, origLine := range page.Content {
        line := Helpers.Trim(origLine)

        if line == ProgramState.Config["blog.excerpt"] {
            break // Exit for loop
        } else {
            excerpt = append(excerpt, line)
        }
    }

    page.Meta["page.excerpt"] = Helpers.Join(excerpt, "\n")
}


/**
 * Name.........: ExpandPage
 * Parameters...: pageInfo (*Page) - the page to expand
 *                ProgarmState (*State.CompilerState) - Compiler state
 * Return.......: error - any error that may have occured
 * Description..: Expands a page into the build directory
 */
func ExpandPage(pageInfo *Page, ProgramState *State.CompilerState) (Errors.Error) {
    page := *pageInfo

    if page.Meta["page.template"] == "" {
        return Errors.NewWarning("No template specified for ", page.File, " it will not be expanded.")
    }

    // Get the contents of the template file, this is a place to start
    contents, err := FileSystem.ReadFile(ProgramState.Template(page.Meta["page.template"] + ".html"))
    if err.HasError() {
        return err
    }

    // Expand the template with the file contents and stuff
    ProgramState.Meta.Push(page.Meta)
    ProgramState.CurrentPage = page
    ExpandContent(&contents, ProgramState)
    ProgramState.Meta.Pop()

    // Write to the output directory
    err = FileSystem.WriteFile(page.OutFile, contents)

    // Perform after file write
    for len(ProgramState.PerformAfterFileWrite) > 0 {
        fn := ProgramState.PerformAfterFileWrite[len(ProgramState.PerformAfterFileWrite) - 1]

        fn(page, ProgramState) // Perform the callback

        ProgramState.PerformAfterFileWrite = ProgramState.PerformAfterFileWrite[:len(ProgramState.PerformAfterFileWrite) - 1]
    }

    return err
}



/**
 * Name.........: ExpandContent
 * Parameters...: content (*[]string) - the content to expand
 *                ProgarmState (*State.CompilerState) - Compiler state
 * Description..: Expands contents for ExpandPage into the build directory
 */
func ExpandContent(content *[]string, ProgramState *State.CompilerState) {
    page := *content

    cmdStack := Stack{}

    partOfCmd := false

    // Loop through lines
    for i, origLine := range page {
        line := Helpers.Trim(origLine)

        // Break up new lines
        lineBreaks := Helpers.Split(line, "\n")
        if len(lineBreaks) > 1 {
            Helpers.Remove(&page, i, i)
            Helpers.Inject(&page, lineBreaks, i)

            ExpandContent(&page, ProgramState)
            break;
        }

        // Check if it an include
        isInclude, fileName := Grammar.IsIncludeStatement(line)
        if isInclude {
            includeContents, err := FileSystem.ReadFile(ProgramState.Include(fileName))
            if !err.HasError() {
                Helpers.Remove(&page, i, i)
                Helpers.Inject(&page, includeContents, i) // Insert the contents of the incldued file

                // Recursively to evaluate stuff in the included file
                ExpandContent(&page, ProgramState)
                break;
            }
        }

        inForEachLoop := false
        if cmdStack.Length() > 0 {
            inForEachLoop = ToCommand(cmdStack.Peek()).Control == "foreach"
        }


        // Only process if statements, foreach loops, and print commands if you are not inside of a foreach loop
        if !inForEachLoop {
            // Check for an inline print statement
            inlinePrints := Grammar.FindInlinePrints(line)

            if len(inlinePrints) > 0 {
                // Evaluate all of the inline commands
                for _, cmd := range inlinePrints {
                    evaluated := Semantics.EvaluatePrintCommand(cmd, ProgramState)

                    page[i] = Helpers.Replace(page[i], cmd, evaluated)
                    line = Helpers.Trim(page[i])
                }

                // Breakup line breaks after evaluating the print command (specifically for {{ content }} or {{ *.excerpt }})
                expanded := Helpers.Split(line, "\n")
                if len(expanded) > 1 {
                    Helpers.Remove(&page, i, i)
                    Helpers.Inject(&page, expanded, i)
                    ExpandContent(&page, ProgramState)
                    break;
                }

            // Look for command that ends something
            } else if Grammar.EndsMultilineCommand(line) {
                whatItEnds := Grammar.WhatDoesEndCommandEnd(line) // Find out what it ends

                // Pull command from stack
                cmd1, _ := cmdStack.Pop()
                cmd := ToCommand(cmd1)
                cmd.EndLine = i
                cmd.State = 0

                // Only process if it ends the command on top of the stack
                if whatItEnds == cmd.Control {
                    toInject := []string{}

                    // Process if statement
                    if cmd.Control == "if" {
                        // Evaluate Conditional
                        if Semantics.IsTrue(cmd.Condition, ProgramState) {
                            toInject = cmd.IfTrue
                        } else {
                            toInject = cmd.IfFalse
                        }

                        // Inject the evaluated if statement where the raw if statement was
                        Helpers.Remove(&page, cmd.StartLine, cmd.EndLine)
                        Helpers.Inject(&page, toInject, cmd.StartLine)

                        ExpandContent(&page, ProgramState)
                        break; // Do not continue processes, the recursive call above will do that
                    }
                } else {
                    cmdStack.Push(cmd) // Put back onto the stack if not ending the topmost command
                }

            // Check if the line starts a multiline command (if, foreach)
            } else if Grammar.StartsMultilineCommand(line) {
                cmd := Semantics.GetCommand(line)
                cmd.StartLine = i

                cmdStack.Push(cmd)
                partOfCmd = true // Do not add this line to the command ifTrue/ifFalse
            }
        } else { // End if !inForEachLoop
            if Grammar.EndsMultilineCommand(line) {
                whatItEnds := Grammar.WhatDoesEndCommandEnd(line) // Find out what it ends

                // Pull command from stack
                cmd1, _ := cmdStack.Pop()
                cmd := ToCommand(cmd1)
                cmd.EndLine = i
                cmd.State = 0

                // Only process if it ends the command on top of the stack
                if whatItEnds == cmd.Control {
                    if cmd.Control == "foreach" {
                        // Perform foreach
                        variable, alias := Semantics.ParseForEachCondition(cmd.Condition)

                        foreachResult := []string{}

                        for _, pg := range ProgramState.GetSpecial(variable) {
                            // Rename the meta keys to the alias specified in the foreach loop
                            newMeta := make(map[string]string)
                            for key, val := range pg.Meta {
                                keys := Helpers.Split(key, ".")
                                newKey := keys[0]

                                if len(keys) > 1 {
                                    newKey = alias + "." + Helpers.Join(keys[1:], ".")
                                }
                                newMeta[newKey] = val
                            }

                            thisLoop := Helpers.Copy(cmd.IfTrue)

                            // Push the new meta to the meta stack, and evaluate the for each loop
                            ProgramState.Meta.Push(newMeta)
                            ExpandContent(&thisLoop, ProgramState)
                            ProgramState.Meta.Pop()

                            // Add to the foreach loop results
                            foreachResult = append(thisLoop, foreachResult...)
                        }

                        // Display the foreach loop results
                        Helpers.Remove(&page, cmd.StartLine, cmd.EndLine)
                        Helpers.Inject(&page, foreachResult, cmd.StartLine)

                        ExpandContent(&page, ProgramState)
                        break;
                    }

                } else {
                    cmdStack.Push(cmd) // Put back onto the stack if not ending the topmost command
                }
            }
        }

        // Add line to the command at the top of the stack if there is one
        if cmdStack.Length() > 0 && !partOfCmd {
            cmd1, _ := cmdStack.Pop()
            cmd := ToCommand(cmd1)

            if cmd.State == 0 {
                cmd.IfTrue = append(cmd.IfTrue, origLine)
            } else if cmd.State == 1 {
                cmd.IfFalse = append(cmd.IfFalse, origLine)
            }

            cmdStack.Push(cmd)
        }



/*
        // Check for an inline print statement
        inlinePrints := Grammar.FindInlinePrints(line)

        if len(inlinePrints) > 0 && !inForEachLoop { // Leaves print commands inside for-each loop
            // Evaluate all of the inline commands
            for _, cmd := range inlinePrints {
                evaluated := Semantics.EvaluatePrintCommand(cmd, ProgramState)

                page[i] = Helpers.Replace(page[i], cmd, evaluated)
                line = Helpers.Trim(page[i])
            }
            // Breakup line breaks after evaluating the print command (specifically for {{ content }} or {{ *.excerpt }})
            expanded := Helpers.Split(line, "\n")
            if len(expanded) > 1 {
                Helpers.Remove(&page, i, i)
                Helpers.Inject(&page, expanded, i)
                ExpandContent(&page, ProgramState)
                break;
            }
            inlinePrints = []string{}

        // Check if the line ends a multiline command
        } else if Grammar.EndsMultilineCommand(line) {
            cmd, _ := cmdStack.Pop()
            cmd.EndLine = i
            cmd.State = 0

            toInject := []string{}

            whatItEnds := Grammar.WhatDoesEndCommandEnd(line)

            if whatItEnds == cmd.Control {
                if cmd.Control == "if" && !inForEachLoop {
                    // Evaluate Conditional
                    if Semantics.IsTrue(cmd.Condition, ProgramState) {
                        Helpers.Print("Yellow", cmd.Condition)
                        toInject = cmd.IfTrue
                    } else {
                        toInject = cmd.IfFalse
                    }
                    Helpers.Remove(&page, cmd.StartLine, cmd.EndLine)
                    Helpers.Inject(&page, toInject, cmd.StartLine)

                    ExpandContent(&page, ProgramState)
                    break;

                } else if cmd.Control == "foreach" {
                    // Perform foreach
                    variable, alias := Semantics.ParseForEachCondition(cmd.Condition)

                    foreachResult := []string{}

                    for _, pg := range ProgramState.GetSpecial(variable) {
                        // Rename the meta keys to the alias specified in the foreach loop
                        Helpers.Print("Green", pg.Meta["page.file"])
                        newMeta := make(map[string]string)
                        for key, val := range pg.Meta {
                            keys := Helpers.Split(key, ".")
                            newKey := keys[0]

                            if len(keys) > 1 {
                                newKey = alias + "." + Helpers.Join(keys[1:], ".")
                            }
                            newMeta[newKey] = val
                        }

                        thisLoop := Helpers.Copy(cmd.IfTrue)
                        for _, ln := range thisLoop {
                            Helpers.Print("red", "\t", ln)
                        }
                        // Push the new meta to the meta stack, and evaluate the for each loop
                        ProgramState.Meta.Push(newMeta)
                        ExpandContent(&thisLoop, ProgramState)
                        ProgramState.Meta.Pop()

                        // Add to the foreach loop results
                        foreachResult = append(thisLoop, foreachResult...)
                    }

                    // Display the foreach loop results
                    Helpers.Remove(&page, cmd.StartLine, cmd.EndLine)
                    Helpers.Inject(&page, foreachResult, cmd.StartLine)
                    ExpandContent(&page, ProgramState)
                    break;
                }
            } else {
                cmdStack.Push(cmd)
            }

        // Check for else statement
        } else if Grammar.StartsElseCommand(line) {
            if cmdStack.Length() > 0 {
                cmd, _ := cmdStack.Pop()

                if cmd.Control == "if" {
                    cmd.State = 1 // Change the state
                }

                cmdStack.Push(cmd)
            }
            partOfCmd = true

        // Check if the line starts a multiline command (if, while, etc...)
        } else if Grammar.StartsMultilineCommand(line) {
            if !inForEachLoop {
                cmd := Semantics.GetCommand(line)
                cmd.StartLine = i

                cmdStack.Push(cmd)
                partOfCmd = false
            } else {
                Helpers.Print("Green", "\t", line)
            }
        }

        // Add line to the command at the top of the stack if there is one
        if cmdStack.Length() > 0 && !partOfCmd {
            cmd, _ := cmdStack.Pop()
            if cmd.State == 0 {
                cmd.IfTrue = append(cmd.IfTrue, origLine)
            } else if cmd.State == 1 {
                cmd.IfFalse = append(cmd.IfFalse, origLine)
            }

            cmdStack.Push(cmd)
        }
*/
        partOfCmd = false
    } // End loop through lines

   *content = page // Give back to caller function
}
