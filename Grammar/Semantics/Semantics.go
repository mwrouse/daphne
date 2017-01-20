package Semantics

import (
    "daphne/State"
    "daphne/Helpers"
    "daphne/Grammar"
    "daphne/Grammar/Operators"
    "daphne/DataTypes"
    "daphne/FileSystem"
    "regexp"
)


var variableRegex, _ = regexp.Compile("^([a-z])+(\\.[a-z]+)+$")

/**
 * Evaluates something to true when given the current state
 */
func IsTrue(inp string, ProgramState *State.CompilerState) (bool) {
    inp = Helpers.Trim(Helpers.ToLower(inp))

    if Grammar.IsLogic(inp) {
        return EvaluateLogicalCondition(inp, ProgramState)
    } else if Grammar.IsCondit(inp) {
        // Evaluate equality
        return EvaluateConditional(inp, ProgramState)

    } else if Grammar.IsStringLit(inp) {
        inp = Helpers.StripQuotes(inp)
    }

    if inp != "" {
        // variable
        val := ""
        if ProgramState.Exists(inp) {
            val = ProgramState.Get(inp)
        } else {
            if variableRegex.MatchString(Helpers.Trim(Helpers.ToLower(inp))) {
                return false;
            }
        }

        val = Helpers.Trim(Helpers.ToLower(val))

        return val != "false" && val != "0" && val != ""
    }

    return false
}



/**
 * Evaluates a Ternary Operator
 */
func EvaluateTernary(ternary string, ProgramState *State.CompilerState) (string) {
    ternary = Helpers.Trim(ternary)

    // Detect if it is a ternary operator or not
    isTernary, condition, ifTrue, ifFalse := Grammar.IsTernary(ternary)

    if !isTernary {
        return ternary
    }

    // It is a ternary, check if the true or false are ternarys, and keep going deeper
    isFalseTernary, _, _, _ := Grammar.IsTernary(ifFalse)
    if isFalseTernary {
        ifFalse = EvaluateTernary(ifFalse, ProgramState) // Recursively solve this
    }

    // Now check if the true statement is a ternary
    isTrueTernary, _, _, _ := Grammar.IsTernary(ifTrue)
    if isTrueTernary {
        ifTrue = EvaluateTernary(ifTrue, ProgramState) // Recursively solve
    }

    // Now, check if the condition is a ternary
    isConditionTernary, _, _, _ := Grammar.IsTernary(condition)
    if isConditionTernary {
        condition = EvaluateTernary(condition, ProgramState)
    }

    if IsTrue(condition, ProgramState) {
        return ifTrue
    } else {
        return ifFalse
    }
}


/**
 * Evaluates a Conditional Expression
 */
func EvaluateConditional(condition string, ProgramState *State.CompilerState) (bool) {

    compound, _, _, _ := Grammar.IsLogicalCondition(condition)
    if compound {
        return EvaluateLogicalCondition(condition, ProgramState)
    }

    isConditional, lhs, operator, rhs := Grammar.IsConditional(condition)
    if !isConditional {
        return false
    }

    lhs = Helpers.ToLower(EvaluateVariable(lhs, ProgramState))
    rhs = Helpers.ToLower(EvaluateVariable(rhs, ProgramState))

    switch (operator) {
    case Operators.Comparison.Equal:
        return lhs == rhs

    case Operators.Comparison.NotEqual:
        return lhs != rhs

    case Operators.Comparison.LargerOrEqual:
        return lhs >= rhs

    case Operators.Comparison.SmallerOrEqual:
        return lhs <= rhs

    case Operators.Comparison.Larger:
        return lhs > rhs

    case Operators.Comparison.Smaller:
        return lhs < rhs
    }

    return false
}


/**
 * Evaluates a compound conditional
 */
func EvaluateLogicalCondition(condition string, ProgramState *State.CompilerState) (bool) {
    compound, lhs, operator, rhs := Grammar.IsLogicalCondition(condition)

    if !compound {
        return false
    }

    lhsBool := IsTrue(lhs, ProgramState)
    rhsBool := IsTrue(rhs, ProgramState)

    switch (operator) {
    case "||":
        return lhsBool || rhsBool
    case "&&":
        return lhsBool && rhsBool
    }

    return false
}




/**
 * Evaluates a Variable
 */
func EvaluateVariable(variable string, ProgramState *State.CompilerState) (string) {
    variable = Helpers.Trim(variable)

    if ProgramState.Exists(Helpers.ToLower(variable)) {
        return ProgramState.Get(Helpers.ToLower(variable))
    }

    // Evaluate concatenation if there is any
    tokens := Helpers.Split(variable, "+")
    if len(tokens) > 1 {
        eval := ""
        for _, token := range tokens {
            eval = eval + EvaluateVariable(token, ProgramState)
        }
        variable = eval
    }

    // Remove quotes from string literals
    if Grammar.IsStringLit(variable) {
        variable = Helpers.StripQuotes(variable)
    } else {
        if variableRegex.MatchString(Helpers.Trim(Helpers.ToLower(variable))) {
            return ""
        }
    }

    return variable
}


/**
  * Name.........:
  * Parameters...:
  * Return.......: string - the evaluated print command
  * Description..: Evaluates a print command
  */
func EvaluatePrintCommand(command string, ProgramState *State.CompilerState) (string) {
    result := Helpers.Trim(command[2:len(command) - 2]) // Remove tags

    if result == "" {
        return result
    }

    eval := ""

    // Check if a special function
    if Grammar.IsSpecialFunc(result) {
        // Evaluate as a function
        eval = EvaluateFunction(result, ProgramState)
    } else {
        // Evaluate not as a function
        eval = EvaluateVariable(EvaluateTernary(result, ProgramState), ProgramState)
    }

    return eval
}


/**
  * Name.........: EvaluateFunction
  * Parameters...: line (string) - the function to validate
  *                ProgramState (*State.CompilerState) - The program state
  * Return.......: string
  * Description..: Evaluates a function
  */
func EvaluateFunction(line string, ProgramState *State.CompilerState) (string) {
    eval := ""

    isValid, funcName, funcParam := Grammar.IsSpecialFunction(line)

    // Return the line if it is not a valid function
    if !isValid {
        return line
    }

    // Handle the function
    switch (funcName) {
    case "post_image":
        funcParams := Helpers.Split(funcParam, ",")
        eval = Helpers.Trim(funcParams[0])

        // Register the function to move the image file after parsing
        ProgramState.PerformAfterFileWrite = append(ProgramState.PerformAfterFileWrite, CopyPostImages(funcParams))
    }

    return eval
}


func CopyPostImages(images []string) (State.SpecialFunction) {
    return func (page DataTypes.Page, ProgramState *State.CompilerState) {
        // Copy all of the images
        for _, img := range images {
            dir := Helpers.Split(page.OutFile, "\\")
            dest := Helpers.Join(dir[:len(dir) - 1], "\\") + "\\" + img

            // Copy the image into the path of the final post
            err := FileSystem.CopyFile(ProgramState.Config["compiler.posts_image_dir"] + "\\" + page.GetSlug() + "\\" + img, dest)
            err.Handle()
        }
    }
}


/**
  * Name.........:
  * Parameters...:
  * Return.......:
  * Description..:
  */
func GetCommand(cmd string) (DataTypes.MultilineCommand) {
    cmd = Helpers.Strip(cmd) // Remove everything around the command

    command := DataTypes.MultilineCommand{State:0}

    cmd = Helpers.Trim(cmd[2:len(cmd) - 2]) // Strip tags

    tokens := Helpers.Split(cmd, " ")

    if len(tokens) < 2 {
        return command
    }

    command.Control = tokens[0]
    command.Condition = Helpers.Join(tokens[1:], " ")

    return command
}


/**
  * Name.........:
  * Parameters...:
  * Return.......:
  * Description..:
  */
 func ParseForEachCondition(condition string) (string, string) {
     operands := Helpers.Split(condition, " as ")

     if len(operands) < 2 {
         return "", ""
     }

     return Helpers.Trim(operands[0]), Helpers.Trim(operands[1])
 }
