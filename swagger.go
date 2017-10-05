/**
 * Swagger.go
 */
package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type Swaggerf struct {
	Swagger Swagger
}

func (s *Swaggerf) ParseSwaggerConfig(jsonBytes []byte) {

	json.Unmarshal(jsonBytes, &s.Swagger)

}

// BuildSwagger builds a swagger file
func (s *Swaggerf) BuildSwagger(rootPath string) {

	log.Printf("Building swagger file from path %s", rootPath)
	sio := &Sio{}
	sio.GetAllFilePaths(rootPath)

	allRoutes := map[string][]Route{}
	allModels := map[string]Model{}

	for _, path := range sio.TmpFiles {

		lines, err := readLines(path)
		if err != nil {
			log.Fatal(err)
		}

		routeMap, _ := GetRoutes(lines, path)

		for path, routes := range routeMap {

			if _, ok := allRoutes[path]; !ok {
				allRoutes[path] = []Route{}
			}

			for _, route := range routes {
				allRoutes[path] = append(allRoutes[path], route)
			}
		}

		models, _ := GetModels(lines, path)

		for name, model := range models {
			allModels[name] = model
		}
	}

	s.Swagger.Paths = map[string]map[string]Path{}

	for pathName, routes := range allRoutes {
		s.Swagger.Paths[pathName] = map[string]Path{}
		s.Swagger.Tags = []Tag{}
		for _, route := range routes {
			path := Path{}
			// path.Description = route.Comments[0]
			// path.OperationID = route.Comments[0]
			path.Description = route.Description
			path.OperationID = route.OperationID
			path.Consumes = []string{"application/json"}
			path.Produces = []string{"application/json"}
			path.Parameters = []Parameter{}
			for _, param := range route.Params {
				parameter := Parameter{}
				parameter.In = param.In
				parameter.Required = param.Required
			}
			path.Responses = map[string]PathResponse{}
			for _, response := range route.Responses {
				pr := PathResponse{}
				pr.Description = response.Description

				if len(response.SchemaRef) > 0 && response.SchemaRef != "empty" {
					pr.Schema = PathSchema{}
					if response.SchemaRef[0:2] == "[]" {
						pr.Schema.Type = "array"
						pr.Schema.Items = map[string]string{}
						pr.Schema.Items["$ref"] = "#/definitions/" + response.SchemaRef[2:]
					} else {
						pr.Schema.Ref = "#/definitions/" + response.SchemaRef
					}
				}

				path.Responses[strconv.Itoa(response.ResponseCode)] = pr
			}
			s.Swagger.Paths[pathName][strings.ToLower(route.Verb)] = path
		}
	}

	s.Swagger.Definitions = map[string]ModelDefinition{}

	for _, model := range allModels {

		definition := ModelDefinition{}

		definition.Type = "object"
		definition.Properties = map[string]Property{}

		for _, field := range model.Fields {
			property := Property{}
			property.Type = field.Type

			definition.Properties[field.Name] = property
		}

		s.Swagger.Definitions[model.Name] = definition
	}
}
