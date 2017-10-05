/**
 * Routes
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// GetRoutes returns a map of route arrays indexed by their path
func GetRoutes(lines []string, filePath string) (routes map[string][]Route, err error) {

	routes = map[string][]Route{}

	var symbols []Symbol

	// Find where the symbols live
	symbols, err = GetSymbols(lines, fmt.Sprintf("@%s", TagRoute))

	if err != nil {
		return
	}

	// If no symbols are found, skip
	if len(symbols) == 0 {
		return
	}

	for _, symbol := range symbols {

		route := Route{}
		comments, _, blockEnd := GetCommentBlock(lines, symbol.LineNum)

		// Parse the symbols inside this comment block
		symbolMap := ParseSymbols(comments)

		// Check that the route field exists (only use the first)
		// Note: There may be more than one route tag, but anything after the first is ignored
		if _, ok := symbolMap[TagRoute]; !ok {
			log.Printf("No routes found at path %s", filePath)
			continue
		}

		if len(symbolMap[TagRoute]) < 1 {
			log.Printf("No routes found at path %s", filePath)
			continue
		}

		route, routeErr := ParseRoute(symbolMap[TagRoute][0], blockEnd+1, filePath, comments)

		if routeErr != nil {
			log.Printf("Route Error: %s", routeErr.Error())
		}

		// Return tags
		if _, ok := symbolMap[TagReturn]; ok {
			for _, ret := range symbolMap[TagReturn] {

				response, err := ParseRouteResponse(ret)

				if err != nil {
					continue
				}

				route.Responses = append(route.Responses, response)
			}
		}

		// Param tags
		if _, ok := symbolMap[TagParam]; ok {
			for _, ret := range symbolMap[TagParam] {
				param, err := ParseRouteParam(ret)
				if err != nil {
					continue
				}
				route.Params = append(route.Params, param)
			}
		}

		if _, ok := symbolMap[TagTags]; ok {
			for _, ret := range symbolMap[TagTags] {
				tags, err := ParseRouteTag(ret)
				if err != nil {
					continue
				}

				route.Tags = tags
			}
		}

		if filePath == "../cherry/api/routes/tasks.go" && route.Path == "/tasks/{task_id}" && route.Verb == "PUT" {
			fmt.Printf("%s - %s %s\n", filePath, route.Verb, route.Path)
			fmt.Printf("%v\n", symbolMap)
			fmt.Printf("%v\n", route.Tags)
		}
		if _, ok := routes[route.Path]; !ok {
			routes[route.Path] = []Route{}
		}

		routes[route.Path] = append(routes[route.Path], route)
	}

	return
}

// ParseRoute parses a route from a line
func ParseRoute(line string, lineNum int, filePath string, comments []string) (route Route, err error) {

	// Assume that after the block ends, so the method starts
	route.LineNum = lineNum
	route.FilePath = filePath
	routeParts := strings.Fields(line)

	if len(routeParts) < 2 {
		err = fmt.Errorf("The tag @route is not in the correct format. File: %s; Line: %d", filePath, route.LineNum)
		return
	}

	route.OperationID = routeParts[0]
	route.Verb = routeParts[1]
	route.Path = routeParts[2]

	if len(comments) > 0 {
		route.Description = comments[0]
	}

	return
}

func ParseRouteParam(ret string) (param Param, err error) {

	retParts := strings.Fields(ret)

	retPartLen := len(retParts)

	if retPartLen < 2 {
		// invalid

	}

	param.Name = retParts[0]
	switch {
	case retParts[1][0:3] == GoTypeInt:
		param.Type = SwaggerTypeInt
	case retParts[1] == GoTypeBool:
		param.Type = SwaggerTypeBool
	case retParts[1] == GoTypeString:
		param.Type = SwaggerTypeString
	case len(retParts[1]) >= 5 && retParts[1][0:5] == GoTypeFloat:
		param.Type = SwaggerTypeFloat
	default:
		param.Type = retParts[1]
	}
	param.Required = true

	curIdx := 2
	maxIdx := 4
	for {
		if retPartLen <= curIdx {
			break
		}
		if len(retParts[curIdx]) > 3 && retParts[curIdx][0:len(TagArgTransportPrefix)] == TagArgTransportPrefix {

			ins := []string{
				TransportPath,
				TransportQuery,
				TransportForm,
				TransportHeader,
				TransportBody,
			}

			if !inArray(retParts[curIdx][3:], ins) {
				err = fmt.Errorf("Invalid transport '%s'", retParts[curIdx][len(TagArgTransportPrefix):])
				return
			}

			param.In = retParts[curIdx][3:]
		}

		if len(retParts[curIdx]) == 8 && (retParts[curIdx] == TagArgOptional || retParts[curIdx] == TagArgRequired) {
			param.Required = retParts[curIdx] == TagArgRequired
		}

		curIdx = curIdx + 1
		if maxIdx <= curIdx {
			break
		}
	}

	// add the description
	if retPartLen >= curIdx {
		param.Description = strings.Join(retParts[curIdx:], " ")
	}

	return
}

// ParseRouteResponse parses a route's response tag (@return)
// Example: @return 200 Foo Returns a Foo object
func ParseRouteResponse(ret string) (response Response, err error) {
	retParts := strings.Fields(ret)

	retPartLen := len(retParts)

	responseInt, _ := strconv.Atoi(retParts[0])

	response.ResponseCode = responseInt

	if retPartLen > 1 {
		response.SchemaRef = retParts[1]
	}

	if retPartLen > 2 {
		response.Description = strings.Join(retParts[2:], " ")
	}

	return
}

// ParseRouteTag parses a route tag line
// Example: @tag foo,bar
func ParseRouteTag(ret string) (tags []string, err error) {

	if len(ret) == 0 {
		err = errors.New("Tags cannot be empty")
		return
	}

	if strings.Index(ret, ",") > -1 {
		tags = strings.Split(ret, ",")
	} else {
		tags = []string{
			ret,
		}
	}

	return
}
