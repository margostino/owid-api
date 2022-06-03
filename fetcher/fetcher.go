package fetcher

import (
	"encoding/csv"
	"fmt"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/configuration"
	"github.com/margostino/owid-api/model"
	"github.com/margostino/owid-api/utils"
	"io"
	"net/http"
	"strconv"
)

type Dataset struct {
	Row map[string]*float64
}

type DatasetIndex struct {
	Index map[string]Dataset
}

var datasetCache = make(map[string]DatasetIndex)
var config = configuration.GetConfig()

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

// Fetch TODO: caching or preloading?
func Fetch(queryResolver string, entity string, year int) map[string]*float64 {
	yearAsString := strconv.Itoa(year)
	dataset := utils.ToSnakeCase(queryResolver)
	dataKey := entity + yearAsString

	if value, ok := datasetCache[dataset]; ok {
		if result, ok := value.Index[dataKey]; ok {
			return result.Row
		}
	}

	var results = make(map[string]*float64)

	url := fmt.Sprintf("%s/%s.yml", config.MetadataBaseUrl, dataset)
	metadata := getMetadata(url)
	data, err := fetchCSVFromUrl(metadata.DataBaseUrl)
	common.Check(err)

	var index = make(map[int]string)
	for idx, row := range data {
		if idx == 0 {
			for dataIndex, column := range row[2:] {
				index[dataIndex] = utils.NormalizeName(column)
			}
			continue
		}

		if row[0] == entity && row[1] == yearAsString {
			for rowIndex, value := range row[2:] {
				resultValue, _ := strconv.ParseFloat(value, 8)
				results[index[rowIndex]] = &resultValue
			}

			datasetMapping := make(map[string]Dataset)
			cacheValue := Dataset{
				Row: results,
			}
			datasetMapping[dataKey] = cacheValue
			datasetIndex := DatasetIndex{
				Index: datasetMapping,
			}
			datasetCache[dataset] = datasetIndex
			return results
		}
	}
	return results
}
