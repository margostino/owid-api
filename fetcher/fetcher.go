package fetcher

import (
	"encoding/csv"
	"fmt"
	"github.com/margostino/owid-api/common"
	"github.com/margostino/owid-api/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Dataset struct {
	Row map[string]*float64
}

type DatasetIndex struct {
	Index map[string]Dataset
}

var datasetCache = make(map[string]DatasetIndex)
var indexCache = loadIndex()

func loadIndex() map[string]string {
	log.Println("Command-line arguments: " + os.Args[0])
	// Load index only in server mode
	if len(os.Args) > 0 {
		log.Println("Command-line arguments: " + os.Args[1])
		return nil
	}
	var urls = make(map[string]string)
	metadataUrl := os.Getenv("METADATA_URL")
	datasetUrls := fmt.Sprintf("%s/datasets.yml", metadataUrl)
	resp, err := http.Get(datasetUrls)
	defer resp.Body.Close()
	common.Check(err)
	bodyBytes, err := io.ReadAll(resp.Body)
	common.UnmarshalYamlBytes(bodyBytes, &urls)
	return urls
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

	log.Printf("Query for dataset %s with entity [%s] and year [%s]\n", dataset, entity, yearAsString)

	if value, ok := datasetCache[dataset]; ok {
		if result, ok := value.Index[dataKey]; ok {
			log.Println("Results from cache")
			return result.Row
		}
	}

	var results = make(map[string]*float64)
	url := indexCache[dataset]
	data, err := fetchCSVFromUrl(url)
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
			log.Printf("New entry in cache for dataset %s and entity %s and year %s\n", dataset, entity, yearAsString)
			return results
		}
	}
	return results
}
