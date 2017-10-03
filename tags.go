package main

import (
	"strings"
)

// Tag Constants
const (
	TagDescription = "description"
	TagRoute       = "route"
	TagModel       = "model"
	TagReturn      = "return"
	TagParam       = "param"
)

// Tags is a collection of tagName constants
var Tags []string = []string{
	TagDescription,
	TagRoute,
	TagModel,
	TagReturn,
	TagParam,
}

func GetSymbols(lines []string, symbol string) (symbols []Symbol) {

	currentLine := 0
	numLines := len(lines)

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

func ParseTags(lines []string) (tags map[string][]string) {

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
