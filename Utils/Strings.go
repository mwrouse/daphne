/**
 * This file contains helper functions for strings and string slices
 * Really stupid easy, wrapper functions, for stuff becausse
 * I don't want to have to include "string" in all of my files
 */
package Utils

import (
    "strings"
)


/**
 * Name.........: Replace
 * Parameters...: orig (string) - the string to replace something in
 *                toReplace (string) - the string to find and replace
 *                replaceWith (string) - what to replace toReplace with
 * Return.......: string - the string with the replacements
 * Description..: Will replace a substring with another substring in a string
 */
func Replace(orig string, toReplace string, replaceWith string) (string) {
    return strings.Replace(line, toReplace, replaceWith, -1)
}


/**
 * Name.........: Trim
 * Parameters...: str (string) - the string to trim
 * Return.......: string
 * Description..: Will trim leading and trailing spaces on a string
 */
func Trim(str string) (string) {
    return strings.TrimSpace(str)
}


/**
 * Name.........: ToLower
 * Parameters...: str (string) - the string to convert to lowercase
 * Return.......: string - str in lowercase letters
 * Description..: Will convert a string to completely lowercase
 */
func ToLower(str string) (string) {
    return strings.ToLower(str)
}


/**
 * Name.........: Split
 * Parameters...: str (string) - the string to split
 *                trigger (string) - the character to split the string at
 * Return.......: []string
 * Description..: Will split a string at a specified target
 */
func Split(str string, trigger string) ([]string) {
    return strings.Split(str, trigger)
}


/**
 * Name.........: Join
 * Parameters...: str ([]string) - the string to join
 *                glue (string) - what to put between joints
 * Return.......: string - the result
 * Description..: Will join an array with glue
 */
func Join(str []string, glue string) (string) {
    return strings.Join(str, glue)
}


/**
 * Name.........: Substr
 * Parameters...: str (string) - the string to return substring from
 *                start (int) - the start position
 *                end (int) - the end position
 * Return.......: string
 * Description..: Will return a substring from in a string
 */
func Substr(str string, start int, end int) (string) {
    result := ""
    for i, c := range str {
        // Jump to next if not at start yet
        if i < start {
            continue
        }
        // Stop if past the end
        if i > end {
            break
        }

        // Add to result
        result = result + string(c)
    }

    return result
}


/**
 * Name.........: Copy
 * Parameters...: arr ([]string) - the string slice to copy
 * Return.......: []string
 * Description..: Will return a copy of arr
 */
func Copy(arr []string) ([]string) {
    return append([]string{}, arr...)
}


/**
 * Name.........: Inject
 * Parameters...: arr1 (*[]string) - the first slice
 *                arr2 ([]string) - the slice to place inside arr1
 *                pos (int) - the position to inject at
 * Description..: Will inject arr2 into arr2 at pos
 */
func Inject(arr1 *[]string, arr2 []string, pos int) {
    sl := *arr1
    arr1Start := Copy(sl[:pos])
    arr1End := Copy(sl[pos:])

    sl = append(arr1Start, arr2...)
    sl = append(sl, arr1End...)

   *arr1 = sl
}


/**
 * Name.........: Remove
 * Parameters...: arr (*[]string) - the array to remove from
 *                start (int) - the start position
 *                end (int) - the end position
 * Description..: Will remove all elements in slice between start and end conclusive
 */
func Remove(arr *[]string, start int, end int) {
    sl := *arr

    end = end + 1
    if end > len(sl) {
        end = len(sl)
    }
    if start < 0 {
        start = 0
    }

    arrStart := Copy(sl[:start])
    arrEnd := Copy(sl[end:])

    sli := append(arrStart, arrEnd...)

   *arr = sl
}
