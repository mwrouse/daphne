package Errors


import (
    "daphne/Helpers"
    "os"
)

type ErrorLevel int

// Error Levels
const (
    NoError ErrorLevel = iota
    Warning
    Fatal
)


/**
 * Struct to represent a Daphne Error Message
 */
type Error struct {
    Level ErrorLevel
    Msg string
}

/**
  * Name.........: New
  * Return.......: Error
  * Description..: Generates a new blank error message
  */
func New() (Error) {
    return *(new(Error))
}

func None() (Error) {
    return New()
}


/**
  * Name.........: NewFatal
  * Return.......: Error
  * Description..: Generates a new error message with level fatal
  */
func NewFatal(params ...string) (Error) {
    err := new(Error)
    err.Level = Fatal

    for i := range params {
        err.Msg = err.Msg + params[i]
    }

    return *err
}


/**
  * Name.........: NewWarning
  * Return.......: Error
  * Description..: Generates a new warning error message
  */
func NewWarning(params ...string) (Error) {
    err := new(Error)
    err.Level = Warning

    for i := range params {
        err.Msg = err.Msg + params[i]
    }

    return *err
}

/**
  * Name.........: IsFatal
  * Return.......: bool
  * Description..: determines if an error is fatal or not
  */
func (err Error) IsFatal() (bool) {
    return err.Level == Fatal
}


/**
  * Name.........: HasError
  * Return.......: bool
  * Description..: determines if an error exists
  */
func (err Error) HasError() (bool) {
    return err.Level != NoError
}


/**
  * Name.........: Handle
  * Description..: Handles an error
  */
func (err Error) Handle() {
    if err.HasError() {
        if err.IsFatal() {
            Helpers.Print("Red", err.Msg)
            os.Exit(1)
        } else {
            Helpers.Print("Yellow", "WARNING: ", err.Msg)
        }
    }
}
