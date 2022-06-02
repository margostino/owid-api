package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/margostino/owid-api/fetcher"
	"github.com/margostino/owid-api/utils"

	"github.com/margostino/owid-api/graph/generated"
	"github.com/margostino/owid-api/graph/model"
)

func (r *queryResolver) TimeUseInSwedenStatisticsSweden(ctx context.Context, entity string, year int) (*model.TimeUseInSwedenStatisticsSwedenDataset, error) {
	dataset := utils.ToSnakeCase("TimeUseInSwedenStatisticsSweden")
	fetcher.Fetch(dataset, entity, year)
	var value *float64
	value = new(float64)
	*value = 20.89
	var response = model.TimeUseInSwedenStatisticsSwedenDataset{
		TimeAllocationAverageDayMen: value,
		TimeAllocationWeekdayWomen:  value,
	}
	return &response, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
