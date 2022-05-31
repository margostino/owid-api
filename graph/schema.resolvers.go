package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/margostino/owid-api/graph/generated"
	"github.com/margostino/owid-api/graph/model"
)

func (r *queryResolver) TimeUseInSweden(ctx context.Context, entity string, year int) (*model.TimeUseInSwedenDataset, error) {
	var a = model.TimeUseInSwedenDataset{
		TimeAllocationAverageDayMen: 10.0,
		TimeAllocationWeekdayWomen:  12.0,
	}
	return &a, nil
	//panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
