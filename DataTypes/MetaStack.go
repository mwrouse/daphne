package DataTypes


type MetaStack struct {
    items []map[string]string

}

/**
 * Gets the length of a stack
 */
func (self MetaStack) Length() (int) {
    return len(self.items)
}

/**
 * Peeks at the top of the stack
 */
func (self MetaStack) Peek() (map[string]string) {
    if self.Length() > 0 {
        return self.items[len(self.items) - 1]
    } else {
        return make(map[string]string)
    }
}

/**
 * Pushes an item onto the stack
 */
func (self *MetaStack) Push(item map[string]string) (int) {
    (*self).items = append((*self).items, item)

    return (*self).Length()
}


/**
 * Pops an item from the stack
 */
func (self *MetaStack) Pop() (map[string]string, int) {
    item := make(map[string]string)
    length := (*self).Length()

    if length > 0 {
        item = (*self).items[length - 1]
        (*self).items = (*self).items[:length - 1]
        length = (*self).Length()
    }

    return item, length
}
