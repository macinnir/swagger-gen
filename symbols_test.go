package main

import "testing"

func TestGetSymbols(t *testing.T) {

	lines := []string{
		"// @route foo",
	}

	symbol := " @route"

	symbols, err := GetSymbols(lines, symbol)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	if len(symbols) != 1 {
		t.Error("Symbol collection should have a length of 1")
	}

	if symbols[0].LineNum != 0 {
		t.Error("First member of symbol collection should have a lineNum of `0`")
	}
}

func TestGetSymbols_ShouldThrowError(t *testing.T) {
	lines := []string{}
	symbol := " @route"

	_, err := GetSymbols(lines, symbol)

	if err == nil {
		t.Error("GetSymbols should have returned non-nil error due to empty lines array")
	}
}

func TestGetCommentBlock(t *testing.T) {

	startLine := 2
	commentLen := 3
	blockStart := 1
	blockEnd := 3

	lines := []string{
		"not a comment block",
		"// test comment",
		"// test comment 2",
		"// test comment 3",
		"not a comment block",
	}

	comments, blockStartActual, blockEndActual := GetCommentBlock(lines, startLine)

	commentLenActual := len(comments)

	if commentLenActual != commentLen {
		t.Errorf("Comment length should be %d (actually %d)", commentLen, commentLenActual)
	}

	if blockStartActual != blockStart {
		t.Errorf("Comment blockStart should be %d (actually %d)", blockStart, blockStartActual)
	}

	if blockEnd != blockEndActual {
		t.Errorf("Comment blockEnd should be %d (actually %d)", blockEnd, blockEndActual)
	}
}

func TestInArray(t *testing.T) {
	needle := "foo"
	haystack := []string{
		"foo",
		"bar",
		"baz",
	}

	result := inArray(needle, haystack)

	if result == false {
		t.Error("inArray should have returned true (actually false)")
	}
}

func TestInArray_ShouldReturnFlase(t *testing.T) {
	needle := "foo"
	haystack := []string{
		"bar",
		"baz",
		"quux",
	}

	result := inArray(needle, haystack)

	if result == true {
		t.Error("inArray should have returned false (actually true)")
	}
}

func TestParseTags(t *testing.T) {

	lines := []string{
		"@foo bar",
		"@description This is the description",
		"@baz quux",
		"@route This the route",
		"@model This is a model",
		"@return This is the return",
		"@param This is the param",
	}

	tagMap := ParseTags(lines)

	if len(tagMap) != 5 {
		t.Errorf("TagMap should contain 5 elements (actually %d", len(tagMap))
	}

	mapShouldContainKey(tagMap, "description", t)
	mapShouldContainKey(tagMap, "route", t)
	mapShouldContainKey(tagMap, "model", t)
	mapShouldContainKey(tagMap, "return", t)
	mapShouldContainKey(tagMap, "param", t)
}

func TestParseTags_NoResults(t *testing.T) {

	lines := []string{
		"foo",
		"bar",
		"baz",
	}

	tagMap := ParseTags(lines)

	if len(tagMap) != 0 {
		t.Errorf("TagMap should have returned an empty collection, but instead has a length of %d", len(tagMap))
	}

}

func mapShouldContainKey(testMap map[string][]string, key string, t *testing.T) {
	if _, ok := testMap[key]; !ok {
		t.Errorf("Map should contain key `%s` but does not", key)
	}
}
