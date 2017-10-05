package main

import "testing"

func TestGetRoutes(t *testing.T) {

	lines := []string{
		"// @route foo",
		"type Foo struct {",
		"",
		"  FooProp1 int",
		"  FooProp2 string",
		"  FooProp3 float",
		"  FooProp4WExtra string somethingElse",
		"}",
	}

	routes, err := GetRoutes(lines, "some/file/path")

	if len(routes) != 1 {
		t.Errorf("GetRoutes should have returned %d routes (actually %d)", 1, len(routes))
	}

	if err != nil {
		t.Errorf("GetRoutes should not have returned an error (actually %s)", err.Error())
	}
}

func TestParseRoute(t *testing.T) {

	comments := []string{
		"GetFoo gets foo",
	}
	line := "GetFoo GET /foo"

	route, err := ParseRoute(line, 0, "some/file/path", comments)

	if err != nil {
		t.Errorf("ParseRoute should have returned a `nil` error (actually %s)", err.Error())
	}

	if route.OperationID != "GetFoo" {
		t.Errorf("ParseRoute should have returned a Route with OperationID == '%s' (actually '%s')", "GetFoo", route.OperationID)
	}

	if route.Verb != "GET" {
		t.Errorf("ParseRoute should have returned a Route with Verb == '%s' (actually '%s')", "GET", route.Verb)
	}

	if route.Path != "/foo" {
		t.Errorf("ParseRoute should have returned a Route with Path == '%s' (actually '%s')", "/foo", route.Path)
	}

	if route.Description != comments[0] {
		t.Errorf("ParseRoute should have returned a Route with Description == '%s' (actually '%s')", comments[0], route.Description)
	}

}

func TestParseRouteParam(t *testing.T) {

	line := "foo int in:path optional This is the foo param"
	param, err := ParseRouteParam(line)

	if err != nil {
		t.Errorf("ParseRouteParam should have a nil error (actually %s)", err.Error())
	}

	if param.Type != "integer" {
		t.Errorf("ParseRouteParam should have returned Param with Type == %s (actually %s)", "integer", param.Type)
	}

	if param.In != "path" {
		t.Errorf("ParseRouteParam should have returned Param with In == '%s' (actually '%s')", "route", param.In)
	}

	if param.Required == true {
		t.Error("ParseRouteParam should have returned Param with Required == false (actually true)")
	}

	if param.Description != "This is the foo param" {
		t.Errorf("ParseRouteParam should have returned Param with Description == '%s' (actually '%s')", "This is the foo param", param.Description)
	}
}

func TestParseRouteParam_ShouldReturnError(t *testing.T) {
	line := "foo int in:foo optional This is the foo param"
	_, err := ParseRouteParam(line)

	if err == nil {
		t.Errorf("ParseRouteParam should have returned an error (actually nil)")
	}

	if err.Error() != "Invalid transport 'foo'" {
		t.Errorf("ParseRouteParam should have returned an error with message '%s' (actually '%s')", "Invalid transport 'foo'", err.Error())
	}
}

func TestParseRouteResponse(t *testing.T) {
	line := "200 Foo Returns a Foo object"

	response, err := ParseRouteResponse(line)

	if err != nil {
		t.Errorf("ParseRouteResponse should have returned error == nil (actually '%s')", err.Error())
	}

	if response.ResponseCode != 200 {
		t.Errorf("ParseRouteResponse should have returned a response code of %d (actually %d)", 200, response.ResponseCode)
	}

	if response.Description != "Returns a Foo object" {
		t.Errorf("ParseRouteResponse should have returned Description == '%s' (actually '%s')", "Returns a Foo object", response.Description)
	}
}

func TestParseRouteTag(t *testing.T) {

	line := "foo,bar"

	tags, err := ParseRouteTag(line)

	if err != nil {
		t.Errorf("ParseRouteTag should have returned a nil error (actually '%s'", err.Error())
	}

	if len(tags) != 2 {
		t.Errorf("ParseRouteTag should have retured a string array of %d tags (actually %d)", 2, len(tags))
	}

	if tags[0] != "foo" {
		t.Errorf("ParseRouteTag should have returned its first tag as '%s' (actually '%s'", "foo", tags[0])
	}

	if tags[1] != "bar" {
		t.Errorf("ParseRouteTag should have returned its second tag as '%s' (actually '%s')", "bar", tags[1])
	}
}

func TestParseRouteTag_Single(t *testing.T) {

	line := "foo"

	tags, err := ParseRouteTag(line)

	if err != nil {
		t.Errorf("ParseRouteTag should have returned a nil error (actually '%s'", err.Error())
	}

	if len(tags) != 1 {
		t.Errorf("ParseRouteTag should have retured a string array of %d tags (actually %d)", 2, len(tags))
	}

	if tags[0] != "foo" {
		t.Errorf("ParseRouteTag should have returned its first tag as '%s' (actually '%s'", "foo", tags[0])
	}
}

func TestParseRouteTag_ShouldReturnError(t *testing.T) {

	line := ""

	_, err := ParseRouteTag(line)

	if err == nil {
		t.Errorf("ParseRouteTag should have returned a error of '%s' (actually '%s')", "Tags cannot be empty", err.Error())
	}
}

func shouldHave(f string, m string, a string, t *testing.T) {
	t.Errorf("%s should have returned %s (actually %s)", f, m, a)
}
