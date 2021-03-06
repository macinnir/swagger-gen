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

		// for arrayModelName, modelName := range arrayModels {
		// 	if _, ok := allArrayModels[arrayModelName]; !ok {
		// 		allArrayModels[arrayModelName] = modelName
		// 	}
		// }
	}

	// Definitions (Models)
	s.Swagger.Definitions = map[string]ModelDefinition{}

	for _, model := range allModels {

		definition := ModelDefinition{}

		definition.Type = "object"
		definition.Properties = map[string]Property{}

		for _, field := range model.Fields {
			property := Property{}
			property.Type = field.Type

			if field.Type == "array" {

				property.Type = "array"
				property.Items = map[string]string{}

				// TODO this does not support a simple array of strings
				// if len(field.Ref) > 0 {
				property.Items["$ref"] = "#/definitions/" + field.Ref
				// }

			} else if field.Type == "#object" {
				property.Ref = "#/definitions/" + field.Ref
			}

			definition.Properties[field.Name] = property
		}

		s.Swagger.Definitions[model.Name] = definition
	}

	s.Swagger.Paths = map[string]map[string]Path{}

	for pathName, routes := range allRoutes {
		s.Swagger.Paths[pathName] = map[string]Path{}
		for _, route := range routes {
			path := Path{}
			// path.Description = route.Comments[0]
			// path.OperationID = route.Comments[0]
			path.Description = route.Description
			path.OperationID = route.OperationID
			path.Consumes = []string{"application/json"}
			path.Produces = []string{"application/json"}
			path.Parameters = []Parameter{}
			if len(route.Tags) > 0 {
				path.Tags = route.Tags
			}
			for _, param := range route.Params {

				parameter := Parameter{}
				paramType := param.Type

				parameter.In = param.In
				parameter.Name = param.Name
				parameter.Description = param.Description

				// Check if the return type is a known model
				if _, ok := s.Swagger.Definitions[paramType]; ok {
					parameter.Schema = map[string]string{}
					parameter.Schema["$ref"] = "#/definitions/" + paramType
				} else {
					parameter.Required = param.Required
					parameter.Schema = map[string]string{}
					parameter.Type = paramType
				}
				path.Parameters = append(path.Parameters, parameter)

			}

			// Responses
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
}
