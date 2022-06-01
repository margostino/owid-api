package tooling

import (
	"fmt"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/configuration"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const COMMENT_BLOCK_BEGIN = "\"\"\""
const COMMENT_BLOCK_END = "\"\"\""

func GenerateSchema() {
	a := "6₆"
	if strings.HasSuffix(a, "₆") || strings.Contains(a, "U+208x") {
		println("dsmk")
	}
	config := configuration.GetConfig()
	var metadata Metadata
	var schema string
	var queryVariables = make([]string, 0)

	filepath.Walk(config.MetadataPath, func(path string, info os.FileInfo, err error) error {
		var argumentsList = make([]string, 0)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if strings.HasSuffix(info.Name(), ".yml") {
			var normalizedTypeNameParts = make([]string, 0)
			fmt.Printf("File Name: %s \n", path)
			ymlFile, err := ioutil.ReadFile(path)
			common.Check(err)
			ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
			err = yaml.Unmarshal(ymlFile, &metadata)
			common.Check(err)
			typeNameParts := strings.Split(metadata.Name, "_")
			for _, typeNamePart := range typeNameParts {
				normalizedTypeNameParts = append(normalizedTypeNameParts, strings.Title(strings.ToLower(typeNamePart)))
			}
			typeName := strings.Join(normalizedTypeNameParts, "") + "Dataset"

			description := metadata.Description
			for _, argument := range metadata.Arguments {
				argumentsList = append(argumentsList, fmt.Sprintf("%s: %s!", argument.Name, argument.Type))
			}
			arguments := strings.Join(argumentsList[:], ", ")
			if description != "" {
				sanitizedDescription := strings.ReplaceAll(description, "\n", "")
				sanitizedDescription = strings.ReplaceAll(sanitizedDescription, ". ", ".\n")
				schema += fmt.Sprintf("%s\n%s\n%s\n", COMMENT_BLOCK_BEGIN, sanitizedDescription, COMMENT_BLOCK_END)
			}

			schema += fmt.Sprintf("type %s {\n", typeName)
			for _, variable := range metadata.Variables {
				schema += fmt.Sprintf(" %s: %s\n", variable.Name, variable.Type)
			}
			schema += "}\n\n"
			queryVariables = append(queryVariables, fmt.Sprintf(" %s(%s):%s!\n", metadata.Name, arguments, typeName))
		}

		return nil
	})

	schema += "type Query {\n"
	for _, queryVariable := range queryVariables {
		schema += queryVariable
	}
	schema += "}\n\n"
	schema += "schema {\n query: Query\n}"
	println(schema)

	file, err := os.OpenFile(config.SchemaFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	common.Check(err)
	defer file.Close()
	_, err2 := file.WriteString(schema)
	common.Check(err2)
}
