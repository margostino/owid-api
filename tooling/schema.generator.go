package tooling

import (
	"encoding/json"
	"fmt"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/configuration"
	"github.com/margostino/owid-api/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const COMMENT_BLOCK_BEGIN = "\"\"\""
const COMMENT_BLOCK_END = "\"\"\""

func GenerateSchema() {
	config := configuration.GetConfig()
	var schema string
	var queryVariables = make([]string, 0)
	var dataBaseUrl = os.Getenv("DATA_BASE_URL")
	var dataUrlMapping string

	filepath.Walk(config.MetadataPath, func(path string, info os.FileInfo, err error) error {
		common.Check(err)

		if info.Name() == "datapackage.json" {
			var normalizedTypeNameParts = make([]string, 0)
			var argumentsList = make([]string, 0)
			var dataPackage map[string]interface{}

			// Get file content
			jsonFile, err := ioutil.ReadFile(path)
			common.Check(err)
			json.Unmarshal(jsonFile, &dataPackage)

			datasetName := utils.NormalizeName(dataPackage["name"].(string))
			dataResources := dataPackage["resources"].([]interface{})
			description := dataPackage["description"].(string)
			resourceName := dataPackage["name"].(string)
			dataUrl := fmt.Sprintf("%s/%s/%s.csv", dataBaseUrl, resourceName, resourceName)
			dataUrlMapping += fmt.Sprintf("%s: %s\n", datasetName, dataUrl)

			// First append description above Variable definition
			if description != "" {
				sanitizedDescription := strings.ReplaceAll(description, "\n", "")
				sanitizedDescription = strings.ReplaceAll(sanitizedDescription, ". ", ".\n")
				schema += fmt.Sprintf("%s\n%s\n%s\n", COMMENT_BLOCK_BEGIN, sanitizedDescription, COMMENT_BLOCK_END)
			}

			// Variable Type
			typeNameParts := strings.Split(datasetName, "_")
			for _, typeNamePart := range typeNameParts {
				normalizedTypeNameParts = append(normalizedTypeNameParts, strings.Title(strings.ToLower(typeNamePart)))
			}
			typeName := strings.Join(normalizedTypeNameParts, "") + "Dataset"

			schema += fmt.Sprintf("type %s {\n", typeName)

			// Fields
			for _, resource := range dataResources {
				dataSchema := resource.(map[string]interface{})["schema"]
				fieldsMap := dataSchema.(map[string]interface{})
				fields := fieldsMap["fields"].([]interface{})
				for _, fieldMap := range fields {
					field := fieldMap.(map[string]interface{})
					fieldName := utils.NormalizeName(field["name"].(string))
					fieldType := utils.NormalizeType(field["type"].(string))
					if fieldName == "entity" || fieldName == "year" {
						argumentsList = append(argumentsList, fmt.Sprintf("%s: %s!", fieldName, fieldType))
					} else if len(fieldName) > 0 && fieldName != "\r" {
						schema += fmt.Sprintf(" %s: %s\n", fieldName, fieldType)
					}
				}
			}
			schema += "}\n\n"
			arguments := strings.Join(argumentsList[:], ", ")
			queryVariables = append(queryVariables, fmt.Sprintf(" %s(%s):%s!\n", datasetName, arguments, typeName))
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

	// Save Schema
	save(config.SchemaFile, schema)

	// Save Data URLs mapping
	save(config.UrlMappingFile, dataUrlMapping)

}

func save(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	common.Check(err)
	defer file.Close()
	_, err2 := file.WriteString(data)
	common.Check(err2)
}
