package main

import (
	"context"
	"github.com/google/go-github/v45/github"
	"log"
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
	client := github.NewClient(nil)
	_, datasets, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", "datasets", nil)
	check(err)
	for _, dataset := range datasets {
		normalizer := NameNormalizer{*dataset.Name}
		datasetName := normalizer.Normalize()
		//path := *dataset.Path
		//_, objects, _, err := client.Repositories.GetContents(context.Background(), "owid", "owid-datasets", path, nil)
		//check(err)
		//for _, object := range objects {
		//	println(object.Path)
		//}
		println(datasetName)
	}
	//fmt.Printf("Original: %s\nNormalized: %s\n", *dataset.Name, value)
}

func check(err error) {
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
