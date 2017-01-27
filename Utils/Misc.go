package Utils





/**
 * Name.........:
 * Parameters...:
 * Return.......:
 * Description..:
 */
func ToMap(self interface{}) (map[string]string) {
    v, ok := self.(map[string]string)

    if !ok {
        // Could not convert, return an empty map
        return make(map[string]string)
    }

    return v
}
