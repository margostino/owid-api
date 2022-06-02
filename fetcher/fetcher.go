package fetcher

import (
	"context"
	"encoding/csv"
	"github.com/google/go-github/v45/github"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/model"
	"io"
	"net/http"
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
	//common.UnmarshalYaml(url, &metadata)
	bodyBytes, err := io.ReadAll(resp.Body)
	common.UnmarshalYamlBytes(bodyBytes, &metadata)
	//reader := yaml.NewDecoder(resp.Body)
	//reader.Comma = ','
	//data, err := reader.ReadAll()

	return &metadata
}

func readCSVFromUrl(url string) ([][]string, error) {
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

func Fetch(dataset string, entity string, year int) {
	url := METADATA[dataset]
	//url := "https://raw.githubusercontent.com/owid/owid-datasets/master/datasets/20th%20century%20deaths%20in%20US%20-%20CDC/20th century deaths in US - CDC.csv"
	metadata := getMetadata(url)

	println(metadata)
	//if err != nil {
	//	panic(err)
	//}
	//for idx, row := range data {
	//	// skip header
	//	if idx == 0 {
	//		continue
	//	}
	//
	//	if idx == 6 {
	//		break
	//	}
	//
	//	fmt.Println(row[2])
	//}

}
