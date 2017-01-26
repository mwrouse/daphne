package DataTypes

// The stack
type Stack []interface{}


/**
 * Name.........: Push
 * Parameters...: item (interface{}) - item to push on the stack
 * Return.......: int - new length of the stack
 * Description..: Pushes an item to the top of the stack
 */
func (s *Stack) Push(item interface{}) (int) {
    (*s) = append(*s, item)

    return len(*s)
}


/**
 * Name.........: Pop
 * Return.......: interface{}
 * Description..: Removes the top item from the stack
 */
func (s *Stack) Pop() (interface{}) {
    if len(*s) == 0 {
        return nil
    }

    popped := (*s)[len(*s) - 1]

    (*s) = (*s)[:len(*s) - 1]

    return popped
}


/**
 * Name.........: Peek
 * Return.......: interface{}
 * Description..: Returns the item on top of the stack without removing it
 */
func (s *Stack) Peek() (interface{}) {
    if len(*s) == 0 {
        return nil
    }

    return (*s)[len(*s) - 1]
}


/**
 * Name.........: Length
 * Return.......: int
 * Description..: Returns the length of the stack
 */
func (s *Stack) Length() (int) {
    return len(*s)
}
