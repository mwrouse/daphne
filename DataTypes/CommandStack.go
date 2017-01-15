package DataTypes


type CommandStack struct {
    items []MultilineCommand
}

/**
 * Gets the length of a stack
 */
func (self CommandStack) Length() (int) {
    return len(self.items)
}

/**
 * Peeks at the top of the stack
 */
func (self CommandStack) Peek() (MultilineCommand) {
    if self.Length() > 0 {
        return self.items[len(self.items) - 1]
    } else {
        return MultilineCommand{}
    }
}

/**
 * Pushes an item onto the stack
 */
func (self *CommandStack) Push(item MultilineCommand) (int) {
    (*self).items = append((*self).items, item)

    return (*self).Length()
}


/**
 * Pops an item from the stack
 */
func (self *CommandStack) Pop() (MultilineCommand, int) {
    item := MultilineCommand{}
    length := (*self).Length()

    if length > 0 {
        item = (*self).items[length - 1]
        (*self).items = (*self).items[:length - 1]
        length = (*self).Length()
    }

    return item, length
}
