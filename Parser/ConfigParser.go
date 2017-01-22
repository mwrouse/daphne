package Parser


import (
    "daphne/State"
    "daphne/Grammar"
    "daphne/FileSystem"
    "daphne/Helpers"
    "daphne/Errors"
)




/**
 * Name.........: ParseConfigFile
 * Parameters...: wd (string)- the path the config file will be in
 *                ProgramState (*State.CompilerState) - The state of the compiler
 * Return.......: error - any errors
 * Description..: Parses the configuration file
 */
 func ParseConfigFile(wd string, ProgramState *State.CompilerState) (Errors.Error) {
     file := wd + "\\_config.daphne"

     config := make(map[string]string)

     // Read the file
     contents, err := FileSystem.ReadFile(file)
     if err.HasError() {
         return err
     }

     // Parse the config contents
     config, err = ParseConfig(contents)
     if err.HasError() {
         return err
     }

     // Apply defaults
     ApplyDefaultConfigOptions(config)

     if config["compiler.ignore"] != "" {
         ProgramState.Ignore = append(ProgramState.Ignore, Helpers.Split(config["compiler.ignore"], ",")...)
     }

     // Add to the compiler state
     ProgramState.Config = config
     return Errors.None() // No error
 }


/**
  * Name.........: ApplyDefaultConfigOptions
  * Parameters...: config (map[string]string) - the config to apply defualts to
  * Description..: Adds default options to a configuration if they do not exist
  */
func ApplyDefaultConfigOptions(config map[string]string) {
    defaults := map[string]string{"compiler.source": ".", "compiler.output": "_build", "site.template": "default", "compiler.template_dir": "_templates", "compiler.include_dir": "_includes", "compiler.posts_dir": "_posts", "compiler.posts_image_dir": "_posts\\images", "compiler.drafts_dir": "_posts/_drafts","compiler.tags.meta":"---","compiler.tags.opening": "{%", "compiler.tags.closing": "%}", "compiler.tags.print.opening": "{{", "compiler.tags.print.closing": "}}", "blog.permalink": "/blog/%slug%", "blog.foldericize": "true", "blog.excerpt":"<!-- more -->"}

    for key, val := range defaults {
        if config[key] == "" {
            config[key] = val
        }
    }

    // Add folders to ingore
    toIgnore := []string{config["compiler.include_dir"], config["compiler.template_dir"], config["compiler.output"], config["compiler.posts_image_dir"]}
    if config["compiler.ignore"] != "" {
        config["compiler.ignore"] = config["compiler.ignore"] + ","
    }
    config["compiler.ignore"] = config["compiler.ignore"] + Helpers.Join(toIgnore, ",")


    if config["site.url"] != "" {
        if config["site.url"][len(config["site.url"]) - 1:] != "/" {
            config["site.url"] = config["site.url"] + "/"
        }
    }

    if config["blog.permalink"] != "" {
        if config["blog.permalink"][:1] == "/" {
            config["blog.permalink"] = config["blog.permalink"][1:]
        }
    }

}


/**
  * Name.........: ParseConfig
  * Parameters...: contents ([]string) - contents of a file to parse
  * Return.......: map[string]map[string]string - contents of the config file
  *                error - any errors
  * Description..: Parses an already read config file
  */
func ParseConfig(contents []string) (map[string]string, Errors.Error) {
     config := make(map[string]string)

     lineNum := "0"
     currentSection := ""

     // Loop through all the lines
     for i, line := range contents {
         line = Helpers.Trim(line)
         lineNum = Helpers.ToStr(i)

         switch {
         // Comment
        case Grammar.ConfigRegex.Comment.MatchString(line):
             // Do nothing

         // Section Declaration
         case  Grammar.ConfigRegex.SectionBegin.MatchString(line):
             matches :=  Grammar.ConfigRegex.SectionBegin.FindStringSubmatch(line)

             if len(matches) != 2 {
                 return nil, Errors.NewFatal("Invalid config section declaration: ", line )
             }
             // Turn all section names to lowercase
             matches[1] = Helpers.Trim(Helpers.ToLower(matches[1]))
             if matches[1] == "" {
                 return nil, Errors.NewFatal("Invalid config secction name: ", line)
             }

             if currentSection != "" {
                 currentSection = currentSection + "." + matches[1]
             } else {
                 currentSection = matches[1]
             }

         // Variable
        case Grammar.ConfigRegex.VariableSet.MatchString(line):
             if currentSection == "" {
                 return nil, Errors.NewFatal("Variables must be declared in a section: ", line)
             }

             matches :=  Grammar.ConfigRegex.VariableSet.FindStringSubmatch(line)

             if len(matches) != 3 {
                 return nil, Errors.NewFatal("Invalid config variable declared at line ", lineNum, ": ", line)
             }

             matches[1] = Helpers.Trim(Helpers.ToLower(matches[1]))

             if matches[1] == "" {
                 return nil, Errors.NewFatal("Invalid config variable name on line ", lineNum, ": ", line)
             }

             // Put into configuration
             config[currentSection + "." + matches[1]] = matches[2]


         // Section End
         case  Grammar.ConfigRegex.SectionEnd.MatchString(line):
             if currentSection == "" {
                 return nil, Errors.NewFatal("Section End Found without being in a section ", line)
             }

             matches := Helpers.Split(currentSection, ".")

             if len(matches) > 1 {
                 matches = matches[:len(matches) - 1]
                 currentSection = Helpers.Join(matches, ".")
             } else {
                 currentSection = ""
             }

         }
     }

     if currentSection != "" {
         return nil, Errors.NewFatal("Unexpected end of config file/section")
     }

     return config, Errors.None()
}
