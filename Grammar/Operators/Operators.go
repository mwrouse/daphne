package Operators

type DaphneComparisonOperators struct {
    Equal string
    NotEqual string
    Larger string
    Smaller string
    LargerOrEqual string
    SmallerOrEqual string

    All string
}

var Comparison = DaphneComparisonOperators{Equal:"==",NotEqual:"!=",Larger:">",LargerOrEqual:">=",Smaller:"<",SmallerOrEqual:"<=",All:"(==|!=|>|>=|<|<=)"}



/**
 * Name.........: IsOperator
 * Parameters...: str (string) - string to check
 * Return.......: bool - true or false
 * Description..: Determins if a string is a vlid operator
 */
func IsOperator(inp string) (bool) {
    return IsComparisonOperator(inp)
}


/**
 * Name.........: IsComparisonOperator
 * Parameters...: str (string) - string to check
 * Return.......: bool - true or false
 * Description..: Determines if a string is a valid comparison operator
 */
func IsComparisonOperator(inp string) (bool) {
    return inp == Comparison.Equal || inp == Comparison.NotEqual || inp == Comparison.Larger || inp == Comparison.LargerOrEqual || inp == Comparison.Smaller || inp == Comparison.SmallerOrEqual
}

func IsLogicalOperator(inp string) (bool) {
    return inp == "||" || inp == "&&"
}
