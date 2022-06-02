package fetcher

import (
	"context"
	"encoding/csv"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/model"
	"github.com/margostino/owid-api/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var METADATA = getAllMetadata()

// TODO: for local use local resources
func getAllMetadata() map[string]string {
	var metadataUrls = make(map[string]string)
	client := github.NewClient(nil)
	_, results, _, err := client.Repositories.GetContents(context.Background(), "margostino", "owid-api", "metadata", nil)
	common.Check(err)

	for _, result := range results {
		key := strings.ReplaceAll(result.GetName(), ".yml", "")
		metadataUrls[key] = result.GetDownloadURL()
	}
	return metadataUrls
}

func getMetadata(url string) *model.Metadata {
	var metadata model.Metadata
	resp, err := http.Get(url)
	defer resp.Body.Close()
	common.Check(err)
	bodyBytes, err := io.ReadAll(resp.Body)
	common.UnmarshalYamlBytes(bodyBytes, &metadata)
	return &metadata
}

func fetchCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Fetch TODO: caching?
func Fetch(dataset string, entity string, year string) map[string]*float64 {
	var results = make(map[string]*float64)
	url := METADATA[dataset]
	metadata := getMetadata(url)
	data, err := fetchCSVFromUrl(metadata.DataBaseUrl)
	common.Check(err)

	if err != nil {
		panic(err)
	}

	var index = make(map[int]string)
	for idx, row := range data {
		if idx == 0 {
			for dataIndex, column := range row[2:] {
				index[dataIndex] = utils.NormalizeName(column)
			}
			continue
		}

		if row[0] == entity && row[1] == year {
			for rowIndex, value := range row[2:] {
				resultValue, _ := strconv.ParseFloat(value, 8)
				results[index[rowIndex]] = &resultValue
			}
			return results
		}
	}
	return results
}
