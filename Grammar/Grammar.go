package Grammar

/**
 * Curly Braces are repetition (0+), square brackets represent optional part of a rule
 *
 * Expression := (Term ConditionalOperator Term) | Term
 * Term := integer | decimal | identifier | String literall
 *
 *
 */
import (
    "daphne/Grammar/Operators"
    "daphne/Helpers"
    "daphne/DataTypes"
    "daphne/State"
)


/**
 * Returns true if something is a keyword
 */
func IsKeyword(inp string) (bool) {
    return inp == "if" || inp == "endif" || inp == "else" || inp == "include"
}


/**
 * Determines if something is a string litteral
 */
func IsStringLit(inp string) (bool) {
    if !Helpers.WrappedInQuotes(inp) {
        return false;
    }

    state := 0
    lastLetter := ""
    inp = Helpers.StripQuotes(inp)

    for _, c := range inp {
        letter := string(c)

        switch(state) {
        case 0:
            if letter == "\"" {
                if lastLetter != "\\" {
                    return false
                }
            }
        }
        lastLetter = letter;
    }
    return true;
}



/**
 * Short hand, inline for the IsConditional
 */
func IsCondit(line string) (bool) {
    results, _, _, _ := IsConditional(line)
    return results
}


/**
 * Determines if something is a conditional expression
 */
func IsConditional(line string) (bool, string, string, string) {
    state := 0

    if len(line) < 2 {
        return false, "", "", ""
    }

    lhs := ""
    operator := ""
    rhs := ""

    waitingFor := DataTypes.StringStack{}
    returnToState := 0

    lastLetter := ""

    for _, c := range line {
        letter := string(c)

        switch (state) {
        case 0:
            if letter == "'" || letter == "\""{
                state = 1
                waitingFor.Push(letter)
            } else if letter == "(" {
                state = 1
                waitingFor.Push(")")
            } else if Operators.IsComparisonOperator(lastLetter + letter) || Operators.IsComparisonOperator(letter) {
                lhs = lhs[:len(lhs) - 2]
                operator = lastLetter + letter
                state = 2
            }

            if state != 2 {
                lhs = lhs + letter
            }

            returnToState = 0

        // Waiting to get out of something
        case 1:
            if letter == waitingFor.Peek() && lastLetter != "\\" {
                _, length := waitingFor.Pop()
                if length <= 0 {
                    state = returnToState
                }
            }else if letter == "'" || letter == "\"" {
                waitingFor.Push(letter)
            } else if letter == "(" {
                waitingFor.Push(")")
            }

            if returnToState == 0 {
                lhs = lhs + letter
            } else {
                rhs = rhs + letter
            }

        // Reading the rhs
        case 2:
            if letter == "'" || letter == "\"" {
                state = 1
                waitingFor.Push(letter)
            } else if letter == "(" {
                state = 1
                waitingFor.Push(")")
            }
            rhs = rhs + letter
            returnToState = 2

        }

        lastLetter = letter
    }

    return (state == 2 && waitingFor.Length() == 0), Helpers.Strip(lhs), operator, Helpers.Strip(rhs)
}


/**
 * Short hand, inline for the IsConditional
 */
func IsLogic(line string) (bool) {
    results, _, _, _ := IsLogicalCondition(line)
    return results
}

/**
 * Returns true if a conditional statement is compounded
 */
func IsLogicalCondition(condition string) (bool, string, string, string) {
    state := 0
    if len(condition) < 2 {
        return false, "", "", ""
    }

    lhs := ""
    operator := ""
    rhs := ""

    waitingFor := DataTypes.StringStack{}
    returnToState := 0

    lastLetter := ""

    for _, c := range condition {
        letter := string(c)

        switch (state) {
        case 0:
            if letter == "'" || letter == "\""{
                state = 1
                waitingFor.Push(letter)
            } else if letter == "(" {
                state = 1
                waitingFor.Push(")")
            } else if Operators.IsLogicalOperator(lastLetter + letter) {
                lhs = lhs[:len(lhs) - 2]
                operator = lastLetter + letter
                state = 2
            }

            if state != 2 {
                lhs = lhs + letter
            }

            returnToState = 0

        // Waiting to get out of something
        case 1:
            if letter == waitingFor.Peek() && lastLetter != "\\" {
                _, length := waitingFor.Pop()
                if length <= 0 {
                    state = returnToState
                }
            }else if letter == "'" || letter == "\"" {
                waitingFor.Push(letter)
            } else if letter == "(" {
                waitingFor.Push(")")
            }

            if returnToState == 0 {
                lhs = lhs + letter
            } else {
                rhs = rhs + letter
            }

        // Reading the rhs
        case 2:
            if letter == "'" || letter == "\"" {
                state = 1
                waitingFor.Push(letter)
            } else if letter == "(" {
                state = 1
                waitingFor.Push(")")
            }
            rhs = rhs + letter
            returnToState = 2

        }

        lastLetter = letter
    }

    return (state == 2 && waitingFor.Length() == 0), Helpers.Strip(lhs), operator, Helpers.Strip(rhs)
}




/**
 * Returns true if a ternary statement
 */
func IsTernary(line string) (bool, string, string, string) {

    condition := ""
    ifTrue := ""
    ifFalse := ""

    returnToState := 0

    sm := State.NewStateMachine(0)

    // State 0 - Reading the condition
    sm.AddState(0, func (inp string, lastInp string, stack *DataTypes.DataStack) (int) {
        state := 0

        if (inp == "'" || inp == "\"") && lastInp != "\\" {
            state = 1
            stack.Push(inp)
        } else if inp == "(" && lastInp != "\\" {
            state = 1
            stack.Push(")")
        } else if inp == "?" && lastInp != "\\" {
             state = 2 // Go to state 2
        }

        // Add to the condition if the next state is not 2
        if state != 2 {
            condition = condition + inp
        }

        returnToState = 0
        return state
    })

    // State 1 - Waiting for the end of somthing (parenthesis, quotes)
    sm.AddState(1, func (inp string, lastInp string, stack *DataTypes.DataStack) (int) {
        state := 1

        if inp == stack.Peek() && lastInp != "\\" {
            stack.Pop() // Remove item from stack
            if stack.Length() < 1 {
                state = returnToState // Only return to the previous state when there is no more left
            }
        } else if (inp == "'" || inp == "\"") && lastInp != "\\" {
            stack.Push(inp) // Wait for these to end
        } else if inp == "(" && lastInp != "\\" {
            stack.Push(")") // Wait for parenthesis to end
        }

        // Append to the appropriate place
        if returnToState == 0 {
            condition = condition + inp
        } else if returnToState == 2 {
            ifTrue = ifTrue + inp
        } else if returnToState == 3 {
            ifFalse = ifFalse + inp
        }

        return state
    })

    // State 2 - Reading the true part of the ternary
    sm.AddState(2, func (inp string, lastInp string, stack *DataTypes.DataStack) (int) {
        state := 2
        if (inp == "'" || inp == "\"") && lastInp != "\\" {
            state = 1
            stack.Push(inp)
        } else if inp == "(" && lastInp != "\\" {
            state = 1
            stack.Push(")")
        } else if inp == "?" && lastInp != "\\" {
            state = 1
            stack.Push(":") // Possible ternary operator inside the if part of the ternary
        } else if inp == ":" && lastInp != "\\" {
            state = 3
        }

        // Add to the if true part if not going to state 3
        if state != 3 {
            ifTrue = ifTrue + inp
        }

        returnToState = 2
        return state
    })

    // State 3 - Final State - Reading the false part of the ternary
    sm.AddFinalState(3, func (inp string, lastInp string, stack *DataTypes.DataStack) (int) {
        state := 3
        if (inp == "'" || inp == "\"") && lastInp != "\\" {
            state = 1
            stack.Push(inp)
        } else if inp == "(" && lastInp != "\\" {
            state = 1
            stack.Push(")")
        }

        ifFalse = ifFalse + inp // Add everything to the if false

        returnToState = 3
        return state
    })

    // Run the state machine
    result, err := sm.Run(line)
    err.Handle()


    return result, Helpers.Strip(condition), Helpers.Trim(Helpers.StripParens(ifTrue)), Helpers.Trim(Helpers.StripParens(ifFalse))
}


/**
  * Name.........: FindInlinePrints
  * Parameters...: line (string) - the line to search inline commands for
  * Return.......: []string - an array of all the commands found
  * Description..: Finds inline commands in a line
  */
func FindInlinePrints(line string) ([]string) {
    result := []string{}

    state := 0
    command := ""

    lastLetter := ""

    for _, c := range line {
        letter := string(c)

        switch (state) {
        // Waiting for opening tags
        case 0:
            if letter == "{" && lastLetter == "{" {
                command = "{{"
                state = 1
            }

        // Reading command
        case 1:
            // Add to command
            command = command + letter

            if letter == "}" && lastLetter == "}" {
                state = 0 // This command has finished

                // Add command to the result
                result = append(result, command)
                command = ""
            }
        }
        lastLetter = letter
    }

    return result
}



/**
  * Name.........: IsIncludeStatement
  * Parameters...: line (string) - the line to search
  * Return.......: bool
  *                string - the file to include
  * Description..: Determines if a line is an include statement
  */
func IsIncludeStatement(line string) (bool, string ) {
    line = Helpers.Strip(line)

    if len(line) < 4 {
        return false, "" // Not long enough
    }

    // Make sure it is the only thing on the line
    if line[:2] == "{%" && line[len(line) - 2:] == "%}" {
        // Check if include
        line = Helpers.Trim(line[2:len(line) - 2])

        tokens := Helpers.Split(line, " ")

        if len(tokens) < 2 {
            return false, ""
        }

        if tokens[0] == "include" {
            return true, Helpers.Strip(tokens[1])
        }

        return false, ""
    }

    return false, ""
}




/**
  * Name.........: StartsAMultilineCommand
  * Parameters...: line (string) - the line to search
  * Return.......: boolean
  * Description..: Determines if a line is the start of a multiline command
  */
func StartsMultilineCommand(line string) (bool) {
     line = Helpers.Strip(line) // Strip everything since a command can only be on a line by itself

     if len(line) < 4 {
         return false // Not enoug characters
     }

     // Has a command if the first and last two characters are tags
     if line[:2] == "{%" && line[len(line) - 2:] == "%}" {
         return true
     }
     return false
 }


 /**
   * Name.........: StartsElseCommand
   * Parameters...: line (string) - the line to search
   * Return.......: boolean
   * Description..: Determines if a line is an else statement
   */
func StartsElseCommand(line string) (bool) {
    line = Helpers.Strip(line)

    if len(line) < 10 {
        return false
    }

    if line[:2] == "{%" && line[len(line)-2:] == "%}" {
        line = Helpers.Trim(line[2:len(line) - 2])

        if line == "else" {
            return true
        }
    }

    return false
}


/**
  * Name.........: EndsMultilineCommand
  * Parameters...: line (string) - the line to search
  * Return.......: boolean
  * Description..: Determines if a line is an end statement
  */
func EndsMultilineCommand(line string) (bool) {
    line = Helpers.Strip(line)

    if len(line) < 9 {
        return false // Not long enough
    }

    if line[:2] == "{%" && line[len(line) - 2:] == "%}" {
        line = Helpers.Trim(line[2:len(line) - 2]) // Remove tags

        tokens := Helpers.Split(line, " ")

        if len(tokens) < 2 {
            return false
        }

        if tokens[0] == "end" {
            return true
        }
    }

    return false
}

/**
  * Name.........: EndsMultilineCommand
  * Parameters...: line (string) - the line to search
  * Return.......: boolean
  * Description..: Determines if a line is an end statement
  */
func WhatDoesEndCommandEnd(line string) (string) {
    line = Helpers.Strip(line)

    if len(line) < 9 {
        return "" // Not long enough
    }

    if line[:2] == "{%" && line[len(line) - 2:] == "%}" {
        line = Helpers.Trim(line[2:len(line) - 2]) // Remove tags

        tokens := Helpers.Split(line, " ")

        if len(tokens) < 2 {
            return ""
        }

        if tokens[0] == "end" {
            return Helpers.Trim(tokens[1])
        }
    }

    return ""
}
