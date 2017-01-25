/**
 * This file has helper functions for directories and URLS
 */
package Utils





/**
 * Name.........: URLSafe
 * Parameters...: url (string) - URL to convert
 * Return.......: string
 * Description..: Converts a URL into all URL safe characters
 */
func URLSafe(url string) (string) {
    url = ToLower(Trim(url))
    fUrl := "" // Final URL

    // Remove all non-alphanumeric characters
    for _, c := range url {
        // Replace spaces with dashes
        if c == 32 {
            fUrl = fUrl + "-"

        // If alphanumeric then add it to the final URL
        } else if IsAlphaNum(c) {
            fUrl = fUrl + string(c)

        // Convert backslashes into forward slashes (47 is forward slash, 92 is backslash)
        } else if c == 47 || c == 92 {
            fUrl = fUrl + "/"
        }
    }

    return fUrl
}


/**
 * Name.........: IsInsideDir
 * Parameters...: toCheck (string) - the path to check
 *                toFind (string) - the directory to look for
 * Return.......: bool
 * Description..: Returns true if toCheck is inside of toFind
 */
func IsInsideDir(toCheck string, toFind string) (bool) {
    // Split both at backslashes
    path := Split(toCheck, "\\")
    dir := Split(toFind, "\\")

    // Not inside if shorter length
    if len(path) < len(dir) {
        return false
    }

    // Begin check
    for i, d := range dir {
        d = Trim(d)
        path[i] = Trim(path[i])

        if d != path[i] {
            return false
        }
    }

    return true 
}
