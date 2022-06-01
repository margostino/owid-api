package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/margostino/owid-api/graph/generated"
	"github.com/margostino/owid-api/graph/model"
)

func (r *queryResolver) O20thCenturyDeathsInUsCdc(ctx context.Context, entity string, year int) (*model.O20thCenturyDeathsInUsCdcDataset, error) {
	var value *int
	value = new(int)
	*value = 10
	var response = model.O20thCenturyDeathsInUsCdcDataset{
		AccidentsExclRoadDeaths: value,
		AccidentsTotalDeaths:    value,
	}
	return &response, nil
}

func (r *queryResolver) TimeUseInSweden(ctx context.Context, entity string, year int) (*model.TimeUseInSwedenDataset, error) {
	var value *float64
	value = new(float64)
	*value = 20.89
	var response = model.TimeUseInSwedenDataset{
		TimeAllocationAverageDayMen: value,
		TimeAllocationWeekdayWomen:  value,
	}
	return &response, nil
	//panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
