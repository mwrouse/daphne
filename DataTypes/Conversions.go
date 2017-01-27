package DataTypes



/**
 * Name.........: ToMap
 * Parameters...: self (interface{}) - interface to convert to map
 * Return.......: map[string]string
 * Description..: Converts an interface{] to a map[string]string
 *
 * This is used when pulling stuff from a stack.
 */
func ToMap(self interface{}) (map[string]string) {
    v, ok := self.(map[string]string)

    if !ok {
        // Could not convert, return an empty map
        return make(map[string]string)
    }

    return v
}


/**
 * Name.........: ToCommand
 * Parameters...: self (interface{}) - item to convert to a Multiline Command
 * Return.......: MultilineCommand
 * Description..: Converts an interface{} (most likely from a Stack), to a command
 */
func ToCommand(self interface{}) (MultilineCommand) {
    v, ok := self.(MultilineCommand)

    if !ok {
        // Could not convert, return an empty map
        return *(new(MultilineCommand))
    }

    return v
}
