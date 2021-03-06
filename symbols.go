package main

import (
	"errors"
	"strings"
)

// Tag Constants
const (
	TagDescription        = "description"
	TagRoute              = "route"
	TagModel              = "model"
	TagReturn             = "return"
	TagParam              = "param"
	TagTags               = "tags"
	TagArgRequired        = "required"
	TagArgOptional        = "optional"
	TagArgTransportPrefix = "in:"
	GoTypeInt             = "int"
	SwaggerTypeInt        = "integer"
	GoTypeString          = "string"
	SwaggerTypeString     = "string"
	GoTypeFloat           = "float"
	SwaggerTypeFloat      = "number"
	GoTypeBool            = "bool"
	SwaggerTypeBool       = "boolean"
	TransportPath         = "path"
	TransportQuery        = "query"
	TransportForm         = "form"
	TransportHeader       = "header"
	TransportBody         = "body"
)

// Tags is a collection of tagName constants
var Tags = []string{
	TagDescription,
	TagRoute,
	TagModel,
	TagReturn,
	TagParam,
	TagTags,
}

// GetSymbols returns a collection of symbol objects based on a symbol string
func GetSymbols(lines []string, symbol string) (symbols []Symbol, err error) {

	currentLine := 0
	numLines := len(lines)

	if numLines == 0 {
		err = errors.New("Lines array cannot be empty")
		return
	}

	for {

		if currentLine > (numLines - 1) {
			break
		}

		if strings.Contains(lines[currentLine], symbol) {
			s := Symbol{
				symbol,
				currentLine,
				lines[currentLine],
				[]string{},
				map[string][]string{},
				"",
			}
			symbols = append(symbols, s)
		}

		currentLine = currentLine + 1
	}

	return
}

// GetCommentBlock parses an entire comment block (comments above a function or struct) and returns them as a string array,
// along with the numeric start and end position of the block
func GetCommentBlock(lines []string, startLine int) (comments []string, blockStart int, blockEnd int) {

	currentLine := startLine

	// Go backwards
	for {
		currentLine = currentLine - 1
		// TODO support docblock (`*`)
		if currentLine < 1 || len(lines[currentLine]) < 3 || lines[currentLine][0:3] != "// " {
			// Start of block
			blockStart = currentLine + 1
			break
		}

		// prepend
		comments = append([]string{lines[currentLine][3:]}, comments...)
	}

	// Reset currentLine
	currentLine = startLine
	comments = append(comments, lines[currentLine][3:])

	for {
		currentLine = currentLine + 1
		if len(lines[currentLine]) < 3 || lines[currentLine][0:3] != "// " {
			blockEnd = currentLine - 1
			break
		}
		comments = append(comments, lines[currentLine][3:])
	}

	return
}

func inArray(needle string, haystack []string) bool {
	for _, el := range haystack {
		if el == needle {
			return true
		}
	}

	return false
}

// ParseSymbols looks for comment tags (starts with `@`) and returns what it finds as a multi-dimensional array
func ParseSymbols(lines []string) (tags map[string][]string) {

	tags = map[string][]string{}

	for _, line := range lines {

		if !strings.HasPrefix(line, "@") {
			continue
		}

		lineParts := strings.Fields(line)
		// Remove the `@` symbol
		tagName := lineParts[0][1:]

		if !inArray(tagName, Tags) {
			continue
		}

		if _, ok := tags[tagName]; !ok {
			tags[tagName] = []string{}
		}

		tags[tagName] = append(tags[tagName], strings.Join(lineParts[1:], " "))
	}

	return
}

// TODO scan entire package for files
// TODO reference other models from within a model
// TODO allow for models (@model) and routes (@route) to be identified by their tags -- allow them to live in the same file.
