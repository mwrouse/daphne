/**
 * Just a bunch of misc helper functions
 */
package Helpers

import (
    "strings"
    "strconv"
    "github.com/fatih/color"
)

func Print(params ...string) {
    if len(params) < 2 {
        Print("red", "Print needs at least two arguments")
        return
    }

    clr := params[0]

    msg := ""
    params = params[1:]
    for i := range params {
        msg = msg + params[i]
    }
    headline := color.New(color.FgWhite, color.Bold)
    switch (ToLower(clr)) {
    case "red":
        headline = color.New(color.FgRed, color.Bold)
    case "green":
        headline = color.New(color.FgGreen, color.Bold)
    case "yellow":
        headline = color.New(color.FgYellow, color.Bold)
    case "blue":
        headline = color.New(color.FgBlue, color.Bold)
    case "white":
        headline = color.New(color.FgWhite, color.Bold)
    case "magenta":
        headline = color.New(color.FgMagenta, color.Bold)
    case "cyan":
        headline = color.New(color.FgCyan, color.Bold)
    default:
        headline = color.New(color.FgRed, color.Bold)
        msg = "Invalid Color: " + clr
    }
    headline.Println(msg)

}

/**
 * Replaces a string with another string
 */
func Replace(line string, toReplace string, replaceWith string) (string) {
    return strings.Replace(line, toReplace, replaceWith, -1)
}

/**
 * Converts to a string
 */
func ToStr(i int) (string) {
    return strconv.Itoa(i)
}

func ToString(i int) (string) {
    return ToStr(i)
}


/**
 * Removes trailing spaces
 */
func Trim(str string) (string) {
    return strings.TrimSpace(str)
}


/**
 * To Lowercase
 */
func ToLower(str string) (string) {
    return strings.ToLower(str)
}


/**
 * Returns true if wrapped in quotes
 */
func WrappedInQuotes(inp string) (bool) {
    inp = Trim(inp)

    return len(inp) >= 2 && ((inp[:1] == "\"" && inp[len(inp) - 1:] == "\"") || (inp[:1] == "'" && inp[len(inp) - 1:] == "'"))
}


/**
 * Removes the quotes around a string
 */
func StripQuotes(inp string) (string) {
    inp = Trim(inp)

    if WrappedInQuotes(inp) {
        inp = inp[1:len(inp) - 1]
    }

    return inp
}


/**
 * Returns true if a string is wrapped in parenthesis
 */
func WrappedInParens(inp string) (bool) {
    inp = Trim(inp)

    return len(inp) >= 2 && inp[:1] == "(" && inp[len(inp) - 1:] == ")"
}


/**
 * Removes the parenthesis around a string
 */
func StripParens(inp string) (string) {
    inp = Trim(inp)

    if len(inp) >= 2 && inp[:1] == "(" && inp[len(inp) - 1:] ==  ")" {
        inp = inp[1:len(inp) - 1]
    }
    return inp
}


/**
 * Strip everything
 */
func Strip(inp string) (string) {
    return Trim(StripParens(StripQuotes(inp)))
}


/**
 * Splits a string
 */
func Split(str string, trigger string) ([]string) {
    return strings.Split(str, trigger)
}

/**
 * Joins a string
 */
func Join(str []string, glue string) (string) {
    return strings.Join(str, glue)
}



/**
 * Adds one array string into another at a certain index
 */
func Inject(arr1 *[]string, arr2[]string, i int) {
    sl1 := *arr1
    arr1Start := append([]string{}, sl1[:i]...)
    arr1End := append([]string{}, sl1[i:]...)

    sl1 = append(arr1Start, arr2...)
    sl1 = append(sl1, arr1End...)

    *arr1 = sl1
}


/**
 * Removes a range of indexes in a slice
 */
func Remove(arr *[]string, start int, end int) {
    sl1 := *arr

    end = end + 1;
    if (end > len(sl1)) {
        end = len(sl1)
    }

    arrStart := append([]string{}, sl1[:start]...)
    arrEnd := append([]string{}, sl1[end:]...)

    sl1 = append(arrStart, arrEnd...)
    *arr = sl1
}


/**
 * Copies one slice into another
 */
func Copy(arr []string) ([]string) {
    return append([]string{}, arr...)
}



func URLSafe(toConvert string) (string) {
    old := ToLower(toConvert)
    url := ""

    // Remove non-alphabetical characters
    for _, c := range old {
        if c <= 122 && c >= 97 {
            url = url + string(c)
        } else if c == 32 {
            url = url + "-"
        }
    }

    return url
}



func Substring(str string, start int, end int) (string) {
    result := ""
    for i, c := range str {
        if i >= start && i <= end {
            result = result + string(c)

            if i == end {
                break;
            }
        }
    }
    return result
}
