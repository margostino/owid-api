package tooling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/configuration"
	"github.com/margostino/owid-api/utils"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Variable struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Metadata struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	DataPath    string      `yaml:"data_path"`
	Arguments   []*Variable `yaml:"arguments"`
	Variables   []*Variable `yaml:"variables"`
}

func normalizeName(value string) string {
	normalizer := utils.NameNormalizer{Value: value}
	return normalizer.Normalize()
}

func normalizeType(value string) string {
	normalizer := utils.TypeNormalizer{Value: value}
	return normalizer.Normalize()
}

func GenerateMetadata() {
	config := configuration.GetConfig()
	var accessToken = config.GithubAccessToken
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, datasets, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", "datasets", nil)
	common.Check(err)
	for _, dataset := range datasets {
		var metadata Metadata
		var dataPackage map[string]interface{}
		datasetName := normalizeName(*dataset.Name)
		path := *dataset.Path + "/datapackage.json"
		encodedData, _, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", path, nil)
		common.Check(err)
		content, err := encodedData.GetContent()
		common.Check(err)
		json.Unmarshal([]byte(content), &dataPackage)
		dataResources := dataPackage["resources"].([]interface{})

		metadata.Name = datasetName
		metadata.DataPath = dataResources[0].(map[string]interface{})["path"].(string) // TODO: check index
		metadata.Description = dataPackage["description"].(string)
		metadata.Arguments = make([]*Variable, 0)
		metadata.Variables = make([]*Variable, 0)

		fmt.Printf("\nOriginal: %s  -  Normalized: %s\n", *dataset.Name, datasetName)
		for _, resource := range dataResources {
			schema := resource.(map[string]interface{})["schema"]
			fieldsMap := schema.(map[string]interface{})
			fields := fieldsMap["fields"].([]interface{})
			for _, fieldMap := range fields {
				field := fieldMap.(map[string]interface{})
				fieldName := normalizeName(field["name"].(string))
				fieldType := normalizeType(field["type"].(string))
				variable := Variable{
					Name: fieldName,
					Type: fieldType,
				}
				if fieldName == "entity" || fieldName == "year" {
					metadata.Arguments = append(metadata.Arguments, &variable)
				} else {
					metadata.Variables = append(metadata.Variables, &variable)
				}
				fmt.Printf("FieldName: %s  -  FieldType: %s\n", fieldName, fieldType)
			}
		}

		yamlData, err := yaml.Marshal(&metadata)
		common.Check(err)
		fileName := fmt.Sprintf("%s/%s.yml", config.MetadataPath, metadata.Name)
		err = ioutil.WriteFile(fileName, yamlData, 0644)
		common.Check(err)
		fmt.Printf("\nNew file: %s\n", fileName)
	}
}
