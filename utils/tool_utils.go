package utils

import (
	"context"
	"fmt"

	"graphQlDemo/ent"
	"graphQlDemo/ent/tool"
	"github.com/google/uuid"
)

func UpdateToolRating(ctx context.Context, client *ent.Client, toolID uuid.UUID) error {
	reviews, reviewsErr := client.Tool.
		Query().
		Where(tool.ID(toolID)).
		QueryReviews().
		All(ctx)
	if reviewsErr  != nil {
		return fmt.Errorf("failed to get reviews")
	}

	var sum int
	for _, review := range reviews {
		sum += review.Rating
	}

	count := len(reviews)
	average := 0.0
	if count > 0 {
		average = float64(sum) / float64(count)
	}

	_, saveErr := client.Tool.
		UpdateOneID(toolID).
		SetAverageRating(average).
		SetRatingCount(count).
		Save(ctx)
	if saveErr != nil {
		return fmt.Errorf("failed to update tool rating")
	}

	return nil
}
