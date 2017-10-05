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

	init := flag.Bool("init", false, "Initialize the swagger-meta.json file with default values")
	sourceDir := flag.String("s", ".", "The root of the source code you want swagger-gen to scan and build a swagger spec from. Defaults to current directory")
	outDir := flag.String("o", ".", "The path to the directory where the generated swagger file will be output to. Defaults to current directory")
	format := flag.String("f", "json", "Output format. json | yaml. Defaults to json")

	flag.Parse()

	if _, dirErr := os.Stat(path.Dir(*sourceDir)); os.IsNotExist(dirErr) {
		log.Fatal(dirErr)
	}

	swaggerMetaPath := path.Join(*sourceDir, "swagger-meta.json")

	if *init == true {
		swagger := Swagger{}
		swagger.Swagger = "2.0"
		swagger.Info = SwaggerInfo{}
		swagger.Info.Description = "My API Description"
		swagger.Info.TermsOfService = "TOS"
		swagger.Info.Title = "My api title"
		swagger.Info.Version = "0.1.0"
		swagger.Info.Contact = Contact{
			"example@email.com",
		}

		swagger.Info.License = License{
			"Apache 2.0",
			"http://www.apache.org/licenses/LICENSE-2.0.html",
		}

		swagger.Host = "myhost.com"
		swagger.BasePath = "/v1"
		swagger.Tags = []Tag{}
		swagger.Schemes = []string{
			"http",
		}
		toJSON(swagger, swaggerMetaPath)
		log.Printf("Swagger meta file generated at path %s", swaggerMetaPath)
		os.Exit(0)
	}

	jsonBytes, jsonBytesErr := ReadJSONToBytes(swaggerMetaPath)

	if jsonBytesErr != nil {
		log.Fatal(jsonBytesErr)
	}

	swaggerf := Swaggerf{}

	swaggerf.ParseSwaggerConfig(jsonBytes)

	// Build the swagger object
	swaggerf.BuildSwagger(*sourceDir)
	var outFile string
	switch *format {
	case "json":
		outFile = "swagger.json"
		outPath := path.Join(*outDir, outFile)
		toJSON(swaggerf.Swagger, outPath)
	case "yaml":
		outFile = "swagger.yaml"
		outPath := path.Join(*outDir, outFile)
		toYAML(swaggerf.Swagger, outPath)
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
