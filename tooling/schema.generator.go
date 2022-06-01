package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Variable struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Metadata struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Arguments   []*Variable `yaml:"arguments"`
	Variables   []*Variable `yaml:"variables"`
}

func main() {
	var metadata Metadata
	var schema string
	//INDENTATION := "    "
	COMMENT_BLOCK_BEGIN := "\"\"\""
	COMMENT_BLOCK_END := "\"\"\""
	//metadataFile := "/Users/martin.dagostino/workspace/margostino/owid-api/graph/metadata/o20th_century_deaths_in_us_cdc.yml"
	metadataFolder := "/Users/martin.dagostino/workspace/margostino/owid-api/graph/metadata"
	schemaFile := "/Users/martin.dagostino/workspace/margostino/owid-api/graph/schema.graphql"
	var queryVariables = make([]string, 0)

	filepath.Walk(metadataFolder, func(path string, info os.FileInfo, err error) error {
		var argumentsList = make([]string, 0)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if strings.HasSuffix(info.Name(), ".yml") {
			var normalizedTypeNameParts = make([]string, 0)
			fmt.Printf("File Name: %s \n", path)
			ymlFile, err := ioutil.ReadFile(path)
			check(err)
			ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
			err = yaml.Unmarshal(ymlFile, &metadata)
			check(err)

			//typeName := metadata.Name + "Dataset"
			typeNameParts := strings.Split(metadata.Name, "_")
			for _, typeNamePart := range typeNameParts {
				normalizedTypeNameParts = append(normalizedTypeNameParts, strings.Title(strings.ToLower(typeNamePart)))
			}
			typeName := strings.Join(normalizedTypeNameParts, "") + "Dataset"

			description := metadata.Description
			for _, argument := range metadata.Arguments {
				argumentsList = append(argumentsList, fmt.Sprintf("%s: %s", argument.Name, argument.Type))
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

	file, err := os.Create(schemaFile)
	check(err)
	defer file.Close()
	_, err2 := file.WriteString(schema)
	check(err2)
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
