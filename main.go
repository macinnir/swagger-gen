package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

func main() {

	sourceDir := flag.String("src", ".", "The root of the source code you want swagger-gen to scan and build a swagger spec from. Defaults to current directory")
	outDir := flag.String("out", ".", "The path to the directory where the generated swagger file will be output to. Defaults to current directory")
	format := flag.String("format", "json", "Output format. json | yaml. Defaults to json")

	flag.Parse()

	if _, dirErr := os.Stat(path.Dir(*sourceDir)); os.IsNotExist(dirErr) {
		log.Fatal(dirErr)
	}

	// Build the swagger object
	swagger := BuildSwagger(*sourceDir)
	var outFile string
	switch *format {
	case "json":
		outFile = "swagger.json"
		outPath := path.Join(*outDir, outFile)
		toJSON(swagger, outPath)
	case "yaml":
		outFile = "swagger.yaml"
		outPath := path.Join(*outDir, outFile)
		toYAML(swagger, outPath)
	default:
		log.Fatal("Invalid output format. Should be `json` or `yaml`")
	}
}

func toYAML(swagger Swagger, outFile string) {
	data, _ := yaml.Marshal(&swagger)
	file, _ := os.OpenFile(
		outFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	_, outFileErr := file.Write(data)
	if outFileErr != nil {
		log.Fatal(outFileErr)
	}
}

func toJSON(swagger Swagger, outFile string) {

	// log.Printf("Outfile Extension: %s", path.Ext(outFile))
	// log.Printf("Outfile base: %s", path.Base(outFile))
	// log.Printf("Outfile dir: %s", path.Dir(outFile))

	if _, dirErr := os.Stat(path.Dir(outFile)); os.IsNotExist(dirErr) {
		log.Fatal(dirErr)
	}

	data, err := json.MarshalIndent(swagger, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	file, _ := os.OpenFile(
		outFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	log.Printf("Writing swagger definition to %s", outFile)
	_, outFileErr := file.Write(data)
	if outFileErr != nil {
		log.Fatal(outFileErr)
	}
}
