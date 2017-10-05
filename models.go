/**
 * Models
 */
package main

import (
	"errors"
	"log"
	"strings"
)

// GetModels searches `lines` for @model tags and adds them to the swagger object
func GetModels(lines []string, filePath string) (models map[string]Model, err error) {

	models = map[string]Model{}

	var symbols []Symbol

	// Routes
	symbols, err = GetSymbols(lines, "@model ")

	if err != nil {
		return
	}

	// If no symbols are found, skip
	if len(symbols) == 0 {
		err = errors.New("No symbols found")
		return
	}

	totalLineLen := len(lines)

	for _, symbol := range symbols {
		comments, _, endLine := GetCommentBlock(lines, symbol.LineNum)
		tagMap := ParseTags(comments)

		// Assume that after the end line will be the start of the model definition
		currentLine := endLine + 1
		model := Model{}
		if _, ok := tagMap["model"]; !ok {
			log.Printf("No model tag found at filePath %s", filePath)
			continue
		}
		log.Printf("TagMap %v", tagMap)
		model.Name = tagMap["model"][0]

		for {

			if totalLineLen <= currentLine {
				break
			}

			if len(lines[currentLine]) == 0 {
				currentLine = currentLine + 1
				continue
			}

			line := strings.TrimLeft(lines[currentLine], " ")
			if len(line) > 4 && line[0:5] == "type " {
				currentLine = currentLine + 1
				continue
			}

			if strings.Trim(lines[currentLine], " ") == "}" {
				break
			}

			fieldLineParts := strings.Fields(strings.TrimPrefix(lines[currentLine], " "))
			fieldType := ""
			switch {
			case len(fieldLineParts[1]) > 4 && fieldLineParts[1][0:5] == "float":
				fieldType = "number"
			case fieldLineParts[1][0:3] == "int":
				fieldType = "integer"
			default:
				fieldType = "string"
			}

			field := ModelField{
				fieldLineParts[0],
				fieldType,
			}
			currentLine = currentLine + 1
			model.Fields = append(model.Fields, field)
		}
		models[model.Name] = model
	}

	return
}
