package main

import "testing"

func TestGetModels(t *testing.T) {

	lines := []string{
		"// @model foo",
		"type Foo struct {",
		"",
		"  FooProp1 int",
		"  FooProp2 string",
		"  FooProp3 float",
		"  FooProp4WExtra string somethingElse",
		"}",
	}

	models, err := GetModels(lines, "some/file/path")

	if err != nil {
		t.Errorf("GetModels should not have return an error (actually %s)", err.Error())
	}

	if len(models) != 1 {
		t.Errorf("GetModels should have returned model collection of length 1 (actually %d)", len(models))
	}

	if _, ok := models["foo"]; !ok {
		t.Errorf("GetModels should have returned a map with key of `foo`")
	}

	if len(models["foo"].Fields) != 4 {
		t.Errorf("GetModels should have returned a map with its first model having 4 fields (actually %d)", len(models["foo"].Fields))
	}

	if models["foo"].Fields[0].Name != "FooProp1" {
		t.Errorf("GetModels should have returned a map with its first model having a field of `FooProp1` (actually %d)", len(models["foo"].Fields[0].Name))
	}
}

func TestGetModels_ShouldReturnErrorIfLinesEmpty(t *testing.T) {

	lines := []string{}

	_, err := GetModels(lines, "some/file/path")

	if err == nil {
		t.Errorf("GetModels should have returned an error (actually nil)")
	}
}

func TestGetModels_ShouldReturnErrorIfNoSymbolsFound(t *testing.T) {

	lines := []string{
		"type Foo struct {",
		"  FooProp1 int",
		"  FooProp2 string",
		"  FooProp3WExtra string somethingElse",
		"}",
	}

	_, err := GetModels(lines, "some/file/path")

	if err == nil {
		t.Errorf("GetModels should have returned an error (actually nil)")
	}

	errString := "No symbols found"

	if err.Error() != errString {
		t.Errorf("GetModels returned an error but it should have been %s (actually %s)", errString, err.Error())
	}
}
