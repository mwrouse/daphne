package DataTypes


type StringStack struct {
    items []string

}

/**
 * Gets the length of a stack
 */
func (self StringStack) Length() (int) {
    return len(self.items)
}

/**
 * Peeks at the top of the stack
 */
func (self StringStack) Peek() (string) {
    if self.Length() > 0 {
        return self.items[len(self.items) - 1]
    } else {
        return ""
    }
}

/**
 * Pushes an item onto the stack
 */
func (self *StringStack) Push(item string) (int) {
    (*self).items = append((*self).items, item)

    return (*self).Length()
}


/**
 * Pops an item from the stack
 */
func (self *StringStack) Pop() (string, int) {
    item := ""
    length := (*self).Length()

    if length > 0 {
        item = (*self).items[length - 1]
        (*self).items = (*self).items[:length - 1]
        length = (*self).Length()
    }

    return item, length
}
