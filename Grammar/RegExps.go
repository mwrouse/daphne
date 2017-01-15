package Grammar

import (
    "regexp"
    "daphne/Grammar/Operators"
)


type DaphneConfigRegex struct {
    SectionBegin    *regexp.Regexp
    SectionEnd      *regexp.Regexp
    VariableSet     *regexp.Regexp
    Comment         *regexp.Regexp
}


func RegexConstructor() (*DaphneConfigRegex, *DaphneMetaRegex) {
    var config = new(DaphneConfigRegex)
    config.SectionBegin, _     = regexp.Compile("^(?:\")?([A-Za-z]+)?(?:\")?(?:\\s)*:(?:\\s)*{(?:\\s)*$")
    config.SectionEnd, _       = regexp.Compile("^(?:\\s)*}(?:\\s)*(?:,)?(?:\\s)*$")
    config.VariableSet, _      = regexp.Compile("^(?:\\s)*(?:\")?([A-Za-z_]+)(?:\")?(?:\\s)*:(?:\\s)*(.*)?(?:\\s)*(?:,)?$")
    config.Comment, _          = regexp.Compile("^(?:\\s)*#(.*)?$")

    var meta = new(DaphneMetaRegex)
    meta.VariableSet = config.VariableSet

    return config, meta
}

var ConfigRegex, MetaRegex = RegexConstructor()
/*
var ConfigRegex = DaphneConfigRegex{}
ConfigRegex.SectionBegin, _     = regexp.Compile("^(?:\")?([A-Za-z]+)?(?:\")?(?:\\s)*:(?:\\s)*{(?:\\s)*$")
ConfigRegex.SectionEnd, _       = regexp.Compile("^(?:\\s)*}(?:\\s)*(?:,)?(?:\\s)*$")
ConfigRegex.VariableSet, _      = regexp.Compile("^(?:\\s)*(?:\")?([A-Za-z_]+)(?:\")?(?:\\s)*:(?:\\s)*(.*)?(?:\\s)*(?:,)?$")
ConfigRegex.Comment, _          = regexp.Compile("^(?:\\s)*#(.*)?$")
*/

type DaphneMetaRegex struct {
    VariableSet     *regexp.Regexp
}

//var MetaRegex = DaphneMetaRegex{}
//MetaRegex.VariableSet = ConfigRegex.VariableSet


var ConditionalRegex, _     = regexp.Compile("^(?:\\s)*(.*)?(?:\\s)"  + Operators.Comparison.All + "(?:\\s)(.*)?(?:\\s)*$")
var TernaryRegex, _         = regexp.Compile("^(?:\\s)*\\((.*)?\\) \\? (.*)? : (.*)?(?:\\s)*$")
