/**
 * Routes
 */
package main

import (
	"log"
	"strconv"
	"strings"
)

// GetRoutes returns a map of route arrays indexed by their path
func GetRoutes(lines []string, filePath string) (routes map[string][]Route, err error) {

	routes = map[string][]Route{}

	// Routes
	symbols := GetSymbols(lines, "@route ")

	// If no symbols are found, skip
	if len(symbols) == 0 {
		return
	}

	for _, symbol := range symbols {

		route := Route{}
		comments, _, blockEnd := GetCommentBlock(lines, symbol.LineNum)

		tagMap := ParseTags(comments)

		// Assume that after the block ends, so the method starts
		route.LineNum = blockEnd + 1
		route.FilePath = filePath

		// Parse the route field (only use the first)
		// Note: There may be more than one route tag, but anything after the first is ignored
		log.Printf("Route: %s", tagMap)
		if _, ok := tagMap["route"]; !ok {
			log.Printf("No routes found at path %s", filePath)
			continue
		}

		routeParts := strings.Fields(tagMap["route"][0])
		if len(routeParts) < 2 {
			log.Fatalf("The tag @route is not in the correct format. File: %s; Line: %d", filePath, route.LineNum)
		}

		if len(comments) > 0 {
			route.Description = comments[0]
		}
		route.OperationID = routeParts[0]
		route.Verb = routeParts[1]
		route.Path = routeParts[2]

		// Return tags
		if _, ok := tagMap["return"]; ok {
			for _, ret := range tagMap["return"] {

				response, err := parseRouteResponse(ret)

				if err != nil {
					continue
				}

				route.Responses = append(route.Responses, response)
			}
		}

		// Param tags
		if _, ok := tagMap["param"]; ok {
			for _, ret := range tagMap["param"] {

				param, err := parseRouteParam(ret)
				if err != nil {
					continue
				}
				route.Params = append(route.Params, param)
			}
		}

		if _, ok := tagMap["tag"]; ok {
			for _, ret := range tagMap["tag"] {
				tags, err := parseRouteTag(ret)

				if err != nil {
					continue
				}

				route.Tags = tags
			}
		}

		if _, ok := routes[route.Path]; !ok {
			routes[route.Path] = []Route{}
		}

		routes[route.Path] = append(routes[route.Path], route)
	}

	return
}

func parseRouteParam(ret string) (param Param, err error) {

	retParts := strings.Fields(ret)

	retPartLen := len(retParts)

	if retPartLen < 2 {
		// invalid

	}

	param.Name = retParts[0]
	param.Type = retParts[1]
	param.Required = true

	curIdx := 2
	maxIdx := 4
	for {
		if retPartLen <= curIdx {
			break
		}
		if len(retParts[curIdx]) > 3 && retParts[curIdx][0:3] == "in:" {
			in := retParts[curIdx][3:]
			switch in {
			case "path":
				param.In = "path"
			case "query":
				param.In = "query"
			case "form":
				param.In = "form"
			case "header":
				param.In = "header"
			case "body":
				param.In = in
			default:
				param.In = "query"
			}
		}

		if len(retParts[curIdx]) == 8 && (retParts[curIdx] == "optional" || retParts[curIdx] == "required") {
			param.Required = retParts[curIdx] == "required"
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

func parseRouteResponse(ret string) (response Response, err error) {
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

func parseRouteTag(ret string) (tags []string, err error) {
	if strings.Index(ret, ",") > -1 {
		tags = strings.Split(ret, ",")
	} else {
		tags = []string{
			ret,
		}
	}

	return
}
