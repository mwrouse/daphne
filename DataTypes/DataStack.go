package DataTypes



type DataStack []interface{}

/**
  * Name.........: Push
  * Parameters...: item (interface{}) - Item to push to the stack
  * Return.......: int - new length of the stack
  * Description..: Pushes an item to the top of the stack
  */
func (self *DataStack) Push(item interface{}) (int) {
    (*self) = append(*self, item)
    return len(*self)
}


/**
  * Name.........: Pop
  * Return.......: interface{} - item removed from the stack
  * Description..: Removes the top item of the stack
  */
func (self *DataStack) Pop() (interface{}) {
    popped := (*self)[len(*self) - 1]
    (*self) = (*self)[:len(*self) - 1]
    return popped
}


/**
  * Name.........: Peek
  * Return.......: interface{} - the current item on top of the stack
  * Description..: Gets the item from the top of the stack without removing it
  */
func (self *DataStack) Peek() (interface{}) {
    return (*self)[len(*self) - 1]
}


/**
  * Name.........: Length
  * Return.......: int - the current size of the stack
  * Description..: Gets the length of the stack
  */
func (self *DataStack) Length() (int) {
    return len(*self)
}
