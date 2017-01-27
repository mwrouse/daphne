package DataTypes


/**
 * An inline command
 */
type InlineCommand struct {
    Control     string // The control structure (if, include, for)
    Condition   string // The condition of the control

    IfTrue      string // What to do if condition is true
    IfFalse     string // What to do if condition is false

    StartLine   int
    EndLine     int
}


/**
 * A command over multiple lines
 */
type MultilineCommand struct {
    Control     string // The control structure (if, include, for)
    Condition   string // The condition of the control

    IfTrue      []string // What to do if condition is true
    IfFalse     []string // What to do if condition is false

    StartLine   int
    EndLine     int

    State       int // What is the state of the command
}



func ToCommand(self interface{}) (MultilineCommand) {
    v, ok := self.(MultilineCommand)

    if !ok {
        // Could not convert, return an empty map
        return *(new(MultilineCommand))
    }

    return v
}
