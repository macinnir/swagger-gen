/**
 * Structs
 */

package main

type Swagger struct {
	Swagger             string                     `json:"swagger"`
	Info                SwaggerInfo                `json:"info"`
	Host                string                     `json:"host"`
	BasePath            string                     `json:"basePath"`
	Tags                []Tag                      `json:"tags"`
	Schemes             []string                   `json:"schemes"`
	Paths               map[string]map[string]Path `json:"paths"`
	SecurityDefinitions map[string]interface{}     `json:"securityDefinitions,omitempty"`
	Definitions         map[string]ModelDefinition `json:"definitions"`
}

type SwaggerInfo struct {
	Description    string  `json:"description"`
	Title          string  `json:"title"`
	Version        string  `json:"version"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
}

type Path struct {
	Description string                  `json:"description"`
	Consumes    []string                `json:"consumes,omitempty"`
	Produces    []string                `json:"produces,omitempty"`
	OperationID string                  `json:"operationId,omitempty"`
	Parameters  []Parameter             `json:"parameters,omitempty"`
	Responses   map[string]PathResponse `json:"responses"`
}

type Parameter struct {
	In          string            `json:"in"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required"`
	Schema      map[string]string `json:"schema,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Contact struct {
	Email string `json:"email"`
}

// Symbol represents a symbol that is looked up in a file
type Symbol struct {
	SymbolString string
	LineNum      int
	Line         string
	Comments     []string
	Tags         map[string][]string
	Type         string // model | route
}

// SrcFile represents a file with the `.go` extension within the target project
type SrcFile struct {
	Lines   []string
	Symbols []Symbol
}

// type Method struct {
// 	Name     string
// 	FilePath string
// 	LineNum  int64
// 	Comments []string
// }

type Route struct {
	Description string
	FilePath    string
	LineNum     int
	Verb        string
	Path        string
	OperationID string
	Summary     string
	Comments    []string
	Params      []Param
	Responses   []Response
	Tags        []string
}

type Response struct {
	ResponseCode int
	Description  string
	SchemaRef    string // sets `type: "array"` if prefixed with `[]`
}

type PathResponse struct {
	Description string     `json:"description"`
	Schema      PathSchema `json:"schema,omitempty"`
}

type PathSchema struct {
	Type  string            `json:"type,omitempty"`
	Ref   string            `json:"$ref,omitempty"`
	Items map[string]string `json:"items,omitempty"`
}

type Param struct {
	Name        string
	Description string
	Required    bool
	Produces    string
	Type        string
	In          string // query || path
}

type Tag struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	ExternalDocs TagExternalDocs `json:"externalDocs"`
}

type TagExternalDocs struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Model struct {
	FilePath string
	LineNum  int
	Name     string
	Fields   []ModelField
}

type ModelField struct {
	Name string
	Type string
}

type Config struct {
	BaseDir   string
	MainFile  string
	ModelsDir string
	RoutesDir string
}

type ModelDefinition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type    string      `json:"type"`
	Format  string      `json:"format,omitempty"`
	Enum    []string    `json:"enum,omitempty"`
	Default interface{} `json:"default,omitempty"`
}
