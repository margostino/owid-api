package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/margostino/owid-api/fetcher"
	"github.com/margostino/owid-api/graph/generated"
	"github.com/margostino/owid-api/graph/model"
	"github.com/margostino/owid-api/utils"
)

func (r *queryResolver) TimeUseInSwedenStatisticsSweden(ctx context.Context, entity string, year int) (*model.TimeUseInSwedenStatisticsSwedenDataset, error) {
	dataset := utils.ToSnakeCase("TimeUseInSwedenStatisticsSweden")
	results := fetcher.Fetch(dataset, entity, strconv.Itoa(year))
	var response = model.TimeUseInSwedenStatisticsSwedenDataset{
		TimeAllocationAverageDayMen: results["time_allocation_average_day_men"],
		TimeAllocationWeekdayWomen:  results["time_allocation_weekday_women"],
	}
	return &response, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
