/**
 * This file contains helper functions for numbers
 */
package Utils

import (
    "strconv"
)


/**
  * Name.........: ToStr
  * Parameters...: i (int) - the number to convert
  * Return.......: string
  * Description..: Converts a number to a string
  */
func ToStr(i int) (string) {
    return strconv.Itoa(i)
}


/**
 * Name.........: IsAlphaNum
 * Parameters...: c (int) - ascii number to check
 * Return.......: bool
 * Description..: Returns true if c is the ascii value of a number or a letter 
 */
func IsAlphaNum(c int) (bool) {
    return (c <= 122 && c >= 97) || (c <= 57 && c >= 48)
}
