/**
 * Swagger.go
 */
package main

import (
	"log"
	"strconv"
	"strings"
)

// BuildSwagger builds a swagger file
func BuildSwagger(rootPath string) Swagger {

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

	swagger := Swagger{}

	swagger.Swagger = "2.0"
	swagger.Info = SwaggerInfo{}
	swagger.Info.Description = "Some Description"
	swagger.Info.TermsOfService = "TOS"
	swagger.Info.Title = "Some Title"
	swagger.Info.Version = "0.1.0"
	swagger.Info.Contact = Contact{
		"rob.macinnis@gmail.com",
	}
	swagger.Host = "goalerfy.com"
	swagger.BasePath = "/v1"

	swagger.Info.License = License{
		"Apache 2.0",
		"http://www.apache.org/licenses/LICENSE-2.0.html",
	}

	swagger.Schemes = []string{
		"http",
	}

	swagger.Paths = map[string]map[string]Path{}

	for pathName, routes := range allRoutes {
		swagger.Paths[pathName] = map[string]Path{}
		swagger.Tags = []Tag{}
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
			swagger.Paths[pathName][strings.ToLower(route.Verb)] = path
		}
	}

	swagger.Definitions = map[string]ModelDefinition{}

	for _, model := range allModels {

		definition := ModelDefinition{}

		definition.Type = "object"
		definition.Properties = map[string]Property{}

		for _, field := range model.Fields {
			property := Property{}
			property.Type = field.Type

			definition.Properties[field.Name] = property
		}

		swagger.Definitions[model.Name] = definition
	}

	return swagger
}
