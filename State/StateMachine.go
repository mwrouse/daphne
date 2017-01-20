package State

import (
    "daphne/DataTypes"
    "daphne/Errors"
    "daphne/Helpers"
)


/**
 * Callback for a State Machine State
 */
type StateCallback func(string, string, *DataTypes.DataStack)(int)


/**
 * A struct to represent a state machine
 */
type StateMachine struct {
    States      map[int]StateCallback
    FinalStates []int
    QuickEndStates []int
    State       int // The current state
    Stack       DataTypes.DataStack
}


/**
  * Name.........: NewStateMachine
  * Parameters...: initialState (int) - the initial state
  * Return.......: *StateMachine
  * Description..: Constructor for a state machine
  */
func NewStateMachine(initialState int) (*StateMachine) {
    sm := new(StateMachine)
    sm.State = initialState
    sm.States = make(map[int]StateCallback, 0)
    sm.FinalStates = make([]int, 0)
    sm.QuickEndStates = make([]int, 0)
    sm.Stack = make(DataTypes.DataStack, 0)

    return sm
}


/**
  * Name.........: AddState
  * Parameters...: state (int) - the state number
  *                fn (StateCallback) - The state callback
  * Return.......: *StateMachine
  * Description..: Adds a state to a state machine
  */
func (self *StateMachine) AddState(state int, fn StateCallback)  (*StateMachine) {
    self.States[state] = fn

    return self
}


/**
  * Name.........: AddFinalState
  * Parameters...: state (int) - the state number
  *                fn (StateCallback) - the state callback
  * Return.......: *StateMachine
  * Description..: Adds an accept state to the state machine
  */
func (self *StateMachine) AddFinalState(state int, fn StateCallback) (*StateMachine) {
    self.AddState(state, fn)
    self.FinalStates = append(self.FinalStates, state)
    return self
}


/**
  * Name.........: Run
  * Parameters...: inp (string) - the string to determine if it is accepted
  * Return.......: bool
  *                Errors.Errror - any errors
  * Description..: Determines if a string is accepted by the state machine
  */
func (self *StateMachine) Run(inp string) (bool, Errors.Error) {
    if len(self.FinalStates) < 1 {
        return false, Errors.NewWarning("State Machine has no final states")
    }

    lastLetter := ""
    letter := ""
    lastState := self.State


    // Loop through all characters in the string
    for _, c := range inp {
        letter = string(c)
        lastState  = self.State

        if self.States[self.State] == nil {
            // Return true if in a quick exit state
            for _, exitState := range self.QuickEndStates {
                if exitState == self.State {
                    return true, Errors.None()
                }
            }

            // Not in a quick exit state, return an error
            return false, Errors.NewWarning("Tried to access an invalid state, ", Helpers.ToStr(self.State), " Previous state was ", Helpers.ToStr(lastState), " Input is: ", inp)
        } else {
            // Call the callback for the state
            self.State = self.States[self.State](letter, lastLetter, &self.Stack)
        }

        lastLetter = letter
    }

    // Check if in a final state
    for _, final := range self.FinalStates {
        if final == self.State {
            // String was accepted
            if self.Stack.Length() > 0 {
                // Stack not empty, return true and a warning
                return true, Errors.NewWarning("State machine ended in a final state, but the stack was not empty. Input: ", inp)
            } else {
                return true, Errors.None()
            }
        }
    }

    // Not a final state, string was not accepted
    return false, Errors.None()
}
