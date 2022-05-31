package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
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
	var argumentsList = make([]string, 0)
	INDENTATION := "    "
	COMMENT_BLOCK_BEGIN := "\"\"\""
	//COMMENT_BLOCK_END := "\"\"\""

	ymlFile, err := ioutil.ReadFile("/Users/martin.dagostino/workspace/margostino/owid-api/graph/metadata/20th_century_deaths_in_us_cdc.yml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &metadata)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	typeName := metadata.Name + "Dataset"
	description := metadata.Description

	if description != "" {
		schema += INDENTATION + COMMENT_BLOCK_BEGIN
		schema += description
		schema += INDENTATION + COMMENT_BLOCK_BEGIN
	}
	schema += fmt.Sprintf("type %s {", typeName)
	for _, variable := range metadata.Variables {
		schema += fmt.Sprintf("%s: %s\n", variable.Name, variable.Type)
	}
	for _, argument := range metadata.Arguments {
		argumentsList = append(argumentsList, fmt.Sprintf("%s: %s", argument.Name, argument.Type))
	}
	arguments := strings.Join(argumentsList[:], ",")

	schema += "}\n"
	schema += "type Query {\n"
	schema += fmt.Sprintf("%s(%s):%s!\n", metadata.Name, arguments, typeName)
	schema += "}\n"
	schema += "schema {\nquery: Query\n}"
	println(schema)
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
